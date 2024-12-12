{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "auditor";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-XJcjxgozwZR4F9KfTW2PaqgPAJeKvF9r2SaEmbjGCyA=";
  subPackages = ["."];
}
