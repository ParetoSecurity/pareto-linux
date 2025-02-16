{
  pkgs,
  lib,
  config,
  inputs,
  ...
}: let
  flakePackage = import ./package.nix {inherit pkgs lib;};
  upstream = import inputs.upstream {system = pkgs.stdenv.system;};
in {
  packages = [
    pkgs.alejandra
    pkgs.goreleaser
  ];
  languages.nix.enable = true;
  languages.go.enable = true;
  languages.go.package = upstream.go_1_24;

  # https://devenv.sh/tests/
  enterTest = ''
    go mod verify
    go test -coverprofile=cover.out ./...
    coverage=$(go tool cover -func=cover.out | grep total | awk '{print $3}' | tr -d %)
    if [ $(echo "$coverage" | sed 's/\..*//') -lt 30 ]; then
      echo "Error: Test coverage is below 30% at $coverage%"
      exit 1
    fi
    echo "Test coverage: $coverage%"
  '';

  # https://devenv.sh/pre-commit-hooks/
  pre-commit.hooks = {
    alejandra.enable = true;
    gofmt.enable = true;
    govet.enable = true;
    revive.enable = true;
    staticcheck.enable = true;
    packaga-sha = {
      name = "Verify package.nix hash";
      enable = true;
      pass_filenames = false;
      files = "go.(mod|sum)$";
      entry = ''
        output=$(nix run .# -- --help 2>&1)
        specified=$(echo "$output" | grep "specified:" | awk '{print $2}')
        got=$(echo "$output" | grep "got:" | awk '{print $2}')
        echo "Specified: $specified"
        echo "Got: $got"
        if [ "$specified" != "$got" ]; then
          echo "Mismatch detected, updating package.nix hash from $specified to $got"
          sed -i"" -e "s/$specified/$got/g" ./package.nix
        else
          echo "Hashes match; no update required."
        fi
      '';
    };
  };

  # See full reference at https://devenv.sh/reference/options/
}
