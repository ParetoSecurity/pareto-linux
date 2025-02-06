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
    go test -coverprofile=cover.out ./...
    coverage=$(go tool cover -func=cover.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
    if (( $(echo "$coverage < 25" | bc -l) )); then
      echo "Test coverage is below 25s%: $coverage%"
      exit 1
    fi
    echo "Test coverage is $coverage%"
  '';

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks = {
    alejandra.enable = true;
    gofmt.enable = true;
    golangci-lint.enable = true;
    govet.enable = true;
    nix-run = {
      name = "Verify package.nix hash";
      enable = true;
      pass_filenames = false;
      files = "go.(mod|sum)$";
      entry = "nix run .# -- --help";
    };
  };

  # See full reference at https://devenv.sh/reference/options/
}
