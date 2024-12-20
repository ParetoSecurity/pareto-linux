# pareto-linux
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/ParetoSecurity/pareto-linux/badge)](https://scorecard.dev/viewer/?uri=github.com/ParetoSecurity/pareto-linux)
[![Integration Tests](https://github.com/ParetoSecurity/pareto-linux/actions/workflows/build.yml/badge.svg)](https://github.com/ParetoSecurity/pareto-linux/actions/workflows/build.yml)
[![Unit Tests](https://github.com/ParetoSecurity/pareto-linux/actions/workflows/unit.yml/badge.svg)](https://github.com/ParetoSecurity/pareto-linux/actions/workflows/unit.yml)
[![Release](https://github.com/ParetoSecurity/pareto-linux/actions/workflows/release.yml/badge.svg)](https://github.com/ParetoSecurity/pareto-linux/actions/workflows/release.yml)


Automatically audit your Linux machine for basic security hygiene.

## Installation

### Using Debian/Ubuntu/Pop!_OS/RHEL/Fedora/CentOS

See [https://pkg.paretosecurity.com](https://pkg.paretosecurity.com) for install steps.


#### Quick Start

To run a one-time security audit:

```bash
paretosecurity check
```

or with JSON reporter

```bash
paretosecurity check --json
```

### Using Nix

#### Quick Start without installing anything

To run a one-time security audit without installation:

```bash
nix run github:paretosecurity/pareto-linux -- check
```

or if running from local repo with JSON reporter

```bash
nix run . -- check --json
```

This will analyze your system and provide a security report highlighting potential improvements and vulnerabilities.

<details>
<summary>
  
### Install via nix-channel

</summary>

As root run:

```ShellSession
$ sudo nix-channel --add https://github.com/paretosecurity/pareto-linux/archive/main.tar.gz paretosecurity
$ sudo nix-channel --update
```

#### Install module via nix-channel

Then add the following to your `configuration.nix` in the `imports` list:

```nix
{
  imports = [ <paretosecurity/modules/paretosecurity.nix> ];
}
```

#### Install CLI via nix-channel

To install the `paretosecurity` binary:

```nix
{
  environment.systemPackages = [ (pkgs.callPackage <paretosecurity/pkgs/paretosecurity.nix> {}) ];
}
```

</details>

<details>
<summary>

### Install via Flakes

</summary>

#### Install module via Flakes

```nix
{
  inputs.paretosecurity.url = "github:paretosecurity/pareto-linux";
  # optional, not necessary for the module
  #inputs.paretosecurity.inputs.nixpkgs.follows = "nixpkgs";

  outputs = { self, nixpkgs, paretosecurity }: {
    # change `yourhostname` to your actual hostname
    nixosConfigurations.yourhostname = nixpkgs.lib.nixosSystem {
      # change to your system:
      system = "x86_64-linux";
      modules = [
        ./configuration.nix
        paretosecurity.nixosModules.default
      ];
    };
  };
}
```

#### Install CLI via Flakes

Using [NixOS module](https://wiki.nixos.org/wiki/NixOS_modules)
(replace system "x86_64-linux" with your system):

```nix
{
  environment.systemPackages = [ paretosecurity.packages.x86_64-linux.default ];
}
```

e.g. inside your `flake.nix` file:

```nix
{
  inputs.paretosecurity.url = "github:paretosecurity/pareto-linux";
  # ...

  outputs = { self, nixpkgs, paretosecurity }: {
    # change `yourhostname` to your actual hostname
    nixosConfigurations.yourhostname = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        # ...
        {
          environment.systemPackages = [ paretosecurity.packages.${system}.default ];
        }
      ];
    };
  };
}
```

</details>
