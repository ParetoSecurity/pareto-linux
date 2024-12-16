{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    nix-vm-test.url = "github:numtide/nix-vm-test";
  };

  outputs = inputs @ {
    flake-parts,
    nixpkgs,
    nix-vm-test,
    ...
  }:
    flake-parts.lib.mkFlake {inherit inputs;} {
      systems = nixpkgs.lib.systems.flakeExposed;

      perSystem = {
        config,
        pkgs,
        lib,
        ...
      }: let
        flakePackage = import ./package.nix {inherit pkgs lib;};
      in {
        packages.default = flakePackage;

        checks.test-nixos = pkgs.testers.runNixOSTest {
          name = "pareto";
          nodes.machine = {pkgs, ...}: {
            environment.systemPackages = [flakePackage];

            systemd.sockets."pareto-linux" = {
              wantedBy = ["sockets.target"];
              socketConfig = {
                ListenStream = "/var/run/pareto-linux.sock";
                SocketMode = "0666";
              };
            };

            systemd.services."pareto-linux" = {
              requires = ["pareto-linux.socket"];
              after = ["pareto-linux.socket"];
              wantedBy = ["multi-user.target"];
              serviceConfig = {
                ExecStart = ["${flakePackage}/bin/paretosecurity" "helper" "--verbose" "--socket" "/var/run/pareto-linux.sock"];
                User = "root";
                Group = "root";
                StandardInput = "socket";
                Type = "oneshot";
                RemainAfterExit = "no";
                StartLimitInterval = "1s";
                StartLimitBurst = 100;
                ProtectSystem = "full";
                ProtectHome = true;
                StandardOutput = "journal";
                StandardError = "journal";
              };
            };

            systemd.user.services."pareto-linux-hourly" = {
              wantedBy = ["timers.target"];
              serviceConfig = {
                Type = "oneshot";
                ExecStart = "${flakePackage}/bin/paretosecurity check";
                StandardInput = "null";
              };
            };

            systemd.user.timers."pareto-linux-hourly" = {
              wantedBy = ["timers.target"];
              timerConfig = {
                OnCalendar = "hourly";
                Persistent = true;
              };
            };
          };

          interactive.nodes.machine = {...}: {
            services.openssh.enable = true;
            services.openssh.settings = {
              PermitRootLogin = "yes";
              PermitEmptyPasswords = "yes";
            };
            security.pam.services.sshd.allowNullPassword = true;
            virtualisation.forwardPorts = [
              {
                from = "host";
                host.port = 2222;
                guest.port = 22;
              }
            ];
          };

          testScript = builtins.readFile "${toString ./.}/test/integration/nixos.py";
        };

        packages.test-debian =
          (inputs.nix-vm-test.lib.x86_64-linux.debian."13" {
            sharedDirs.packageDir = {
              source = "${toString ./.}/pkg";
              target = "/mnt/package";
            };
            testScript = builtins.readFile "${toString ./.}/test/integration/debian.py";
          })
          .driver;

        packages.test-fedora =
          (inputs.nix-vm-test.lib.x86_64-linux.fedora."40" {
            sharedDirs.packageDir = {
              source = "${toString ./.}/pkg";
              target = "/mnt/package";
            };
            testScript = builtins.readFile "${toString ./.}/test/integration/fedora.py";
          })
          .driver;

        packages.test-ubuntu =
          (inputs.nix-vm-test.lib.x86_64-linux.ubuntu."23_10" {
            sharedDirs.packageDir = {
              source = "${toString ./.}/pkg";
              target = "/mnt/package";
            };
            testScript = builtins.readFile "${toString ./.}/test/integration/ubuntu.py";
          })
          .driver;
      };
    };
}
