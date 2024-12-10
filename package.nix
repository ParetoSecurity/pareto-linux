{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "auditor";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-xp65jHGRlmky4gNAPWDOomPkeSjQcs04R9Behkha72g=";
  subPackages = ["."];
}
