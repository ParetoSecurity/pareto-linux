{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "auditor";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-827WJ+bIdPHVgSay580v8lLmBkQQFX7UKGOPI9qNNPI=";
  subPackages = ["."];
}
