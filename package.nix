{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "auditor";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-d/EItA3H8yiQ+VSAH+ZA6jH2Ojb7OkRv8eUcqpabNwI=";
  subPackages = ["."];
}
