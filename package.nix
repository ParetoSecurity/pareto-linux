{
  pkgs,
  lib,
}:
pkgs.buildGoModule rec {
  pname = "pareto-linux";
  version = "${builtins.hashFile "sha256" "${toString ./go.sum}"}";
  src = ./.;
  vendorHash = "sha256-0HN9SXsioXdfd7p2SFaQ/oT4fy8pzSM34NiGL67g8es=";
  subPackages = ["."];
}
