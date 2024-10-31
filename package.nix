{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "auditor";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-g9dv0EaltGwdK9en1N2PEP0b+VzyCxK+Tu+TIXkkyBs";
  subPackages = ["."];
}
