# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

The aws-sso-config team takes security seriously. If you discover a security vulnerability, please follow these steps:

### Responsible Disclosure

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities via email to: [your-email-here]

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

### What to Include

Please include the requested information listed below (as much as you can provide) to help us better understand the nature and scope of the possible issue:

* Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
* Full paths of source file(s) related to the manifestation of the issue
* The location of the affected source code (tag/branch/commit or direct URL)
* Any special configuration required to reproduce the issue
* Step-by-step instructions to reproduce the issue
* Proof-of-concept or exploit code (if possible)
* Impact of the issue, including how an attacker might exploit the issue

This information will help us triage your report more quickly.

### Response Process

1. **Acknowledgment**: We'll acknowledge receipt of your vulnerability report within 48 hours.

2. **Assessment**: Our team will assess the vulnerability and determine its severity and impact.

3. **Fix Development**: We'll work on developing a fix for the vulnerability.

4. **Testing**: The fix will be thoroughly tested to ensure it addresses the vulnerability without introducing new issues.

5. **Release**: We'll release a patched version and publish a security advisory.

6. **Credit**: If you'd like, we'll acknowledge your responsible disclosure in our release notes.

### Timeline

* **Initial Response**: Within 48 hours
* **Status Update**: Within 7 days with an estimated timeline for resolution
* **Resolution**: We aim to resolve critical vulnerabilities within 30 days

## Security Best Practices

When using aws-sso-config, we recommend:

1. **Keep Updated**: Always use the latest version of aws-sso-config
2. **Secure Configuration**: Store AWS credentials securely and follow AWS security best practices
3. **Environment Variables**: Use environment variables for sensitive configuration instead of hardcoding
4. **File Permissions**: Ensure AWS configuration files have appropriate permissions (600)
5. **Regular Audits**: Regularly audit your AWS configurations and access patterns

## Security Features

aws-sso-config implements several security features:

- Uses AWS SDK v2 with modern security practices
- Supports AWS SSO for secure authentication
- No storage of long-term credentials
- Secure token handling and caching
- Input validation and sanitization

## Vulnerability Disclosure Timeline

We believe in coordinated disclosure and will work with security researchers to:

1. Confirm and analyze reported vulnerabilities
2. Develop and test fixes
3. Release patches in a timely manner
4. Provide credit to researchers (if desired)

## Contact

For security-related questions or concerns, please contact: [your-email-here]

For general questions, please use the [GitHub issues](https://github.com/blairham/aws-sso-config/issues) page.
