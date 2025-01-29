{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "paretosecurity";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-tqupnxl3DVFpxRzmMX0BNKzWSqeACrfyYLq7JQCtTgs=";
  subPackages = ["."];
  postInstall = ''
    mv $out/bin/pareto-core $out/bin/paretosecurity
  '';
}
