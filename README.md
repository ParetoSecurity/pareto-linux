# pareto-linux
Automatically audit your Linux machine for basic security hygiene.

https://github.com/user-attachments/assets/6c6b2bec-947c-41f8-93f3-f5d15e80e6a8


## Installation

### Using Nix

The recommended way to install Pareto Linux is through the Nix package manager:

```bash
nix profile install --accept-flake-config github:paretosecurity/pareto-linux
```

#### Quick Start

To run a one-time security audit without installation:

```bash
nix run github:paretosecurity/pareto-linux -- check
```

or if running from local repo

```bash
nix run . -- check
```

This will analyze your system and provide a security report highlighting potential improvements and vulnerabilities.