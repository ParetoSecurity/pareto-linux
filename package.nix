{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "auditor";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-C3puhJBW7yWgbaAWSQVwR7GdZXh0GVa0mxEEENRB6Qg=";
  subPackages = ["."];
}
