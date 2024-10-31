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
        self',
        inputs',
        pkgs,
        system,
        lib,
        ...
      }: let
        flakePackage = import ./package.nix {inherit pkgs lib;};
      in {
        packages.default = flakePackage;

        checks.test-nixos = pkgs.testers.runNixOSTest {
          name = "pareto";

          nodes.machine = {
            lib,
            pkgs,
            ...
          }: {
            environment.systemPackages = [flakePackage];
          };

          # To SSH into the generated vm:
          # $ devenv build outputs.test-pareto.driverInteractive
          # Copy derivation path on last line of output
          # $ /nix/store/9m2sny6g7k7i0zln7s14wznmc8hfpcz5-nixos-test-driver-pareto --interactive
          # >>> start_all()
          # Now the VM is running, let's SSH to it:
          # ssh root@localhost -p2222

          interactive.nodes.machine = {...}: {
            services.openssh = {
              enable = true;
              settings = {
                PermitRootLogin = "yes";
                PermitEmptyPasswords = "yes";
              };
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

          testScript = builtins.readFile "${toString ./.}/integration/nixos.py";
        };

        packages.test-debian = let
          vmTest = inputs.nix-vm-test.lib.x86_64-linux.debian."13" {
            sharedDirs = {
              packageDir = {
                source = "${./.}";
                target = "/mnt/package";
              };
            };
            testScript = builtins.readFile "${toString ./.}/integration/debian.py";
          };
        in
          vmTest.driver;

        packages.test-fedora = let
          vmTest = inputs.nix-vm-test.lib.x86_64-linux.fedora."40" {
            sharedDirs = {
              packageDir = {
                source = "${./.}";
                target = "/mnt/package";
              };
            };
            testScript = builtins.readFile "${toString ./.}/integration/fedora.py";
          };
        in
          vmTest.driver;

        packages.test-ubuntu = let
          vmTest = inputs.nix-vm-test.lib.x86_64-linux.ubuntu."23_10" {
            sharedDirs = {
              packageDir = {
                source = "${./.}";
                target = "/mnt/package";
              };
            };
            testScript = builtins.readFile "${toString ./.}/integration/ubuntu.py";
          };
        in
          vmTest.driver;
      };
    };
}
