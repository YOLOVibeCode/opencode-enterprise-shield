# Homebrew Formula for Enterprise Shield
class EnterpriseShield < Formula
  desc "Enterprise-grade security plugin for OpenCode AI assistant"
  homepage "https://github.com/yourorg/opencode-enterprise-shield"
  version "1.0.0"
  
  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/yourorg/opencode-enterprise-shield/releases/download/v1.0.0/enterprise-shield-v1.0.0-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_ARM64"
    else
      url "https://github.com/yourorg/opencode-enterprise-shield/releases/download/v1.0.0/enterprise-shield-v1.0.0-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_AMD64"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/yourorg/opencode-enterprise-shield/releases/download/v1.0.0/enterprise-shield-v1.0.0-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_ARM64"
    else
      url "https://github.com/yourorg/opencode-enterprise-shield/releases/download/v1.0.0/enterprise-shield-v1.0.0-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_AMD64"
    end
  end

  def install
    bin.install "enterprise-shield"
    
    # Install default config to share directory
    (share/"enterprise-shield").install "config/default.yaml"
    
    # Create OpenCode directories
    opencode_dir = "#{Dir.home}/.opencode"
    plugins_dir = "#{opencode_dir}/plugins"
    config_dir = "#{opencode_dir}/config"
    
    mkdir_p plugins_dir
    mkdir_p config_dir
    
    # Symlink binary to OpenCode plugins directory
    ln_sf bin/"enterprise-shield", "#{plugins_dir}/enterprise-shield"
    
    # Install default config if it doesn't exist
    config_file = "#{config_dir}/enterprise-shield.yaml"
    unless File.exist?(config_file)
      cp "#{share}/enterprise-shield/default.yaml", config_file
    end
  end

  def post_install
    opencode_dir = "#{Dir.home}/.opencode"
    
    # Create logs directory
    mkdir_p "#{opencode_dir}/logs/enterprise-shield"
    
    ohai "Enterprise Shield installed successfully!"
    puts ""
    puts "Configuration: #{opencode_dir}/config/enterprise-shield.yaml"
    puts "Logs: #{opencode_dir}/logs/enterprise-shield/"
    puts ""
    puts "Quick Start:"
    puts "  enterprise-shield version         # Check version"
    puts "  enterprise-shield scan <content>  # Scan for violations"
    puts "  enterprise-shield init            # Reinitialize config"
    puts ""
    puts "Documentation: https://github.com/yourorg/opencode-enterprise-shield"
  end

  test do
    system "#{bin}/enterprise-shield", "version"
  end
end

