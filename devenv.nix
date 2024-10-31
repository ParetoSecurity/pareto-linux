{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: let
  flakePackage = import ./package.nix {inherit pkgs lib;};
in {
  packages = [
    pkgs.alejandra
    pkgs.goreleaser
  ];
  languages.nix.enable = true;
  languages.go.enable = true;

  # https://devenv.sh/tests/
  enterTest = ''
    go mod verify
    go test ./...
    go build .
  '';

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks = {
    alejandra.enable = true;
    gofmt.enable = true;
    golangci-lint.enable = true;
    govet.enable = true;
  };

  # See full reference at https://devenv.sh/reference/options/
}
