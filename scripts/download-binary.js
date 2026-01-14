#!/usr/bin/env node
// NPM post-install script to download the appropriate binary

const https = require('https');
const fs = require('fs');
const path = require('path');
const { promisify } = require('util');
const { pipeline } = require('stream');
const zlib = require('zlib');
const tar = require('tar');

const streamPipeline = promisify(pipeline);

const REPO = 'YOLOVibeCode/opencode-enterprise-shield';
const VERSION = require('../package.json').version;

// Detect platform and architecture
function getPlatform() {
  const platform = {
    darwin: 'darwin',
    linux: 'linux',
    win32: 'windows'
  }[process.platform];

  const arch = {
    x64: 'amd64',
    arm64: 'arm64'
  }[process.arch];

  if (!platform || !arch) {
    throw new Error(`Unsupported platform: ${process.platform}-${process.arch}`);
  }

  return { platform, arch };
}

// Download file
async function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Follow redirect
        return download(response.headers.location, dest).then(resolve).catch(reject);
      }

      if (response.statusCode !== 200) {
        reject(new Error(`Download failed: ${response.statusCode}`));
        return;
      }

      response.pipe(file);
      file.on('finish', () => {
        file.close();
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => {});
      reject(err);
    });
  });
}

async function main() {
  console.log('🛡️  Installing Enterprise Shield binary...');

  const { platform, arch } = getPlatform();
  console.log(`Platform: ${platform}-${arch}`);

  const binaryName = platform === 'windows' ? 'enterprise-shield.exe' : 'enterprise-shield';
  const archiveName = `enterprise-shield-v${VERSION}-${platform}-${arch}`;
  const archiveExt = platform === 'windows' ? '.zip' : '.tar.gz';
  const archiveFile = `${archiveName}${archiveExt}`;

  const downloadUrl = `https://github.com/${REPO}/releases/download/v${VERSION}/${archiveFile}`;

  console.log(`Downloading from: ${downloadUrl}`);

  // Create bin directory
  const binDir = path.join(__dirname, '..', 'bin');
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const tempFile = path.join(binDir, archiveFile);
  const binaryPath = path.join(binDir, binaryName);

  try {
    // Download archive
    await download(downloadUrl, tempFile);
    console.log('✓ Downloaded');

    // Extract
    if (platform === 'windows') {
      // For Windows, we'd need to extract .zip
      const AdmZip = require('adm-zip');
      const zip = new AdmZip(tempFile);
      zip.extractAllTo(binDir, true);
    } else {
      // Extract tar.gz
      await streamPipeline(
        fs.createReadStream(tempFile),
        zlib.createGunzip(),
        tar.extract({
          cwd: binDir,
          strip: 0
        })
      );
    }
    console.log('✓ Extracted');

    // Make executable
    if (platform !== 'windows') {
      fs.chmodSync(binaryPath, 0o755);
    }

    // Clean up
    fs.unlinkSync(tempFile);

    console.log('✓ Installation complete!');
    console.log(`Binary installed at: ${binaryPath}`);
  } catch (error) {
    console.error('❌ Installation failed:', error.message);
    console.error('');
    console.error('You can manually download the binary from:');
    console.error(`  ${downloadUrl}`);
    process.exit(1);
  }
}

main().catch(console.error);

