class Aja < Formula
  desc "Deploy applications with managed dependencies in seconds, not hours"
  homepage "https://deployaja.id"
  version "1.0.0"

  if OS.mac?
    if Hardware::CPU.arm?
      url "https://github.com/deployaja/deployaja-cli/releases/download/v#{version}/aja-darwin-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_ARM64"
    else
      url "https://github.com/deployaja/deployaja-cli/releases/download/v#{version}/aja-darwin-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_AMD64"
    end
  elsif OS.linux?
    if Hardware::CPU.arm?
      url "https://github.com/deployaja/deployaja-cli/releases/download/v#{version}/aja-linux-arm64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_ARM64"
    else
      url "https://github.com/deployaja/deployaja-cli/releases/download/v#{version}/aja-linux-amd64.tar.gz"
      sha256 "PLACEHOLDER_SHA256_LINUX_AMD64"
    end
  end

  def install
    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "aja-darwin-arm64" => "aja"
      else
        bin.install "aja-darwin-amd64" => "aja"
      end
    elsif OS.linux?
      if Hardware::CPU.arm?
        bin.install "aja-linux-arm64" => "aja"
      else
        bin.install "aja-linux-amd64" => "aja"
      end
    end
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/aja version")
  end
end 