// TypeScript/JavaScript wrapper for Enterprise Shield Go plugin
// This allows OpenCode to interface with the Go binary via stdio

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');

class EnterpriseShieldPlugin {
  constructor() {
    this.binaryPath = this.findBinary();
    this.process = null;
    this.requestId = 0;
    this.pendingRequests = new Map();
  }

  findBinary() {
    const pluginDir = path.join(os.homedir(), '.opencode', 'plugins');
    const binaryName = process.platform === 'win32' 
      ? 'enterprise-shield.exe' 
      : 'enterprise-shield';
    return path.join(pluginDir, binaryName);
  }

  async start() {
    if (this.process) {
      return; // Already started
    }

    this.process = spawn(this.binaryPath, ['serve'], {
      stdio: ['pipe', 'pipe', 'inherit']
    });

    this.process.stdout.on('data', (data) => {
      this.handleResponse(data.toString());
    });

    this.process.on('error', (error) => {
      console.error('Enterprise Shield error:', error);
    });

    this.process.on('exit', (code) => {
      console.log(`Enterprise Shield exited with code ${code}`);
      this.process = null;
    });

    // Wait a bit for the process to start
    await new Promise(resolve => setTimeout(resolve, 100));
  }

  async stop() {
    if (this.process) {
      this.process.kill();
      this.process = null;
    }
  }

  async call(method, params) {
    await this.start();

    const id = ++this.requestId;
    const request = {
      jsonrpc: '2.0',
      id,
      method,
      params
    };

    return new Promise((resolve, reject) => {
      this.pendingRequests.set(id, { resolve, reject });
      this.process.stdin.write(JSON.stringify(request) + '\n');

      // Timeout after 30 seconds
      setTimeout(() => {
        if (this.pendingRequests.has(id)) {
          this.pendingRequests.delete(id);
          reject(new Error('Request timeout'));
        }
      }, 30000);
    });
  }

  handleResponse(data) {
    const lines = data.split('\n').filter(line => line.trim());
    
    for (const line of lines) {
      try {
        const response = JSON.parse(line);
        
        if (response.id && this.pendingRequests.has(response.id)) {
          const { resolve, reject } = this.pendingRequests.get(response.id);
          this.pendingRequests.delete(response.id);

          if (response.error) {
            reject(new Error(response.error.message));
          } else {
            resolve(response.result);
          }
        }
      } catch (error) {
        console.error('Failed to parse response:', error);
      }
    }
  }

  // Hook implementations
  async beforeRequest({ userId, content, provider }) {
    return this.call('processRequest', {
      userId,
      content,
      provider
    });
  }

  async afterResponse({ content, sessionId }) {
    return this.call('processResponse', {
      content,
      sessionId
    });
  }

  async scan({ content }) {
    return this.call('scanContent', {
      content
    });
  }
}

// OpenCode plugin export
module.exports = {
  name: 'enterprise-shield',
  version: '1.0.0',
  
  async activate() {
    this.plugin = new EnterpriseShieldPlugin();
    console.log('Enterprise Shield plugin activated');
  },

  async deactivate() {
    if (this.plugin) {
      await this.plugin.stop();
    }
    console.log('Enterprise Shield plugin deactivated');
  },

  hooks: {
    async beforeRequest(context) {
      const result = await this.plugin.beforeRequest({
        userId: context.user?.id || 'anonymous',
        content: context.request.content,
        provider: context.provider
      });

      if (result.blocked) {
        throw new Error(`Request blocked: ${result.blockReason}`);
      }

      // Modify the request content
      context.request.content = result.content;
      
      // Store session ID for desanitization
      context.meta = context.meta || {};
      context.meta.enterpriseShield = {
        sessionId: result.sessionId,
        wasSanitized: result.wasSanitized
      };

      return context;
    },

    async afterResponse(context) {
      const shield = context.meta?.enterpriseShield;
      
      if (!shield || !shield.wasSanitized) {
        return context; // Nothing to desanitize
      }

      const result = await this.plugin.afterResponse({
        content: context.response.content,
        sessionId: shield.sessionId
      });

      // Restore original values
      context.response.content = result.desanitizedContent;

      return context;
    }
  }
};

