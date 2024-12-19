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
        self,
        system,
        ...
      }: let
        flakePackage = import ./package.nix {inherit pkgs lib;};
      in {
        packages.default = flakePackage;

        checks.test-nixos = import ./test/integration/nixos.nix {
          inherit pkgs system flakePackage;
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
