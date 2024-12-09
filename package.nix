{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "auditor";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-LX/WvkKJ4M7FmtckTdXXWyDUbsxbOhnwSSh56lhWFzk=";
  subPackages = ["."];
}
