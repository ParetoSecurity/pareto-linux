{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-Cq1Ia+GV+bbFfRPlmpn0ITW/lAYLPqyf1ESw6Km83wQ=";
  subPackages = ["."];
  postInstall = ''
    mv $out/bin/pareto-linux $out/bin/paretosecurity
  '';
}
