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

  scripts.buildIcon.exec = ''
    # Build the icon for the appropriate platform.
    echo "Generating platform icon..."
    go get github.com/cratonica/2goarray
    cat "assets/icon_black.png" | go run github.com/cratonica/2goarray IconBlack shared > shared/icon_black_unix.go
    cat "assets/icon_white.png" | go run github.com/cratonica/2goarray IconWhite shared > shared/icon_white_unix.go
    go mod tidy
  '';

  # https://devenv.sh/tests/
  enterTest = ''
    go mod verify
    go test ./...
    go build .
    goreleaser check
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
