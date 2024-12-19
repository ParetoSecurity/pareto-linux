{pkgs ? import <nixpkgs> {}}: {
  paretosecurity = pkgs.callPackage ./package.nix {};
}
