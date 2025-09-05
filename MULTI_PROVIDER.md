# Multi-Provider Configuration Guide

GoDNS now supports configuring multiple DNS providers simultaneously, allowing you to manage domains across different DNS services from a single GoDNS instance.

## Features

- **Multi-Provider Support**: Configure multiple DNS providers (Cloudflare, DNSPod, DigitalOcean, etc.) in one config file
- **Per-Domain Provider**: Specify different providers for different domains  
- **Full Backward Compatibility**: Existing single-provider configurations continue to work unchanged
- **Mixed Configuration**: Combine legacy global provider with new per-domain providers
- **Provider-Specific Credentials**: Each provider can have its own authentication credentials

## Configuration Options

### Option 1: Multi-Provider Configuration

Configure multiple providers with per-domain assignment:

```json
{
  "providers": {
    "cloudflare": {
      "email": "user@example.com",
      "password": "your-cloudflare-api-token"
    },
    "dnspod": {
      "login_token": "your-dnspod-token"
    },
    "digitalocean": {
      "password": "your-digitalocean-api-token"
    }
  },
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "api"],
      "provider": "cloudflare"
    },
    {
      "domain_name": "mysite.net",
      "sub_domains": ["mail", "ftp"], 
      "provider": "dnspod"
    }
  ]
}
```

### Option 2: Legacy Single-Provider (Backward Compatible)

Existing configurations continue to work unchanged:

```json
{
  "provider": "DNSPod",
  "login_token": "your-dnspod-token",
  "domains": [
    {
      "domain_name": "example.com",
      "sub_domains": ["www", "test"]
    }
  ]
}
```

### Option 3: Mixed Configuration

Combine global provider with specific per-domain providers:

```json
{
  "provider": "DNSPod", 
  "login_token": "your-dnspod-token",
  "providers": {
    "cloudflare": {
      "email": "user@example.com",
      "password": "your-cloudflare-api-token"
    }
  },
  "domains": [
    {
      "domain_name": "oldsite.com",
      "sub_domains": ["www", "mail"]
    },
    {
      "domain_name": "newsite.com", 
      "sub_domains": ["www", "api"],
      "provider": "cloudflare"
    }
  ]
}
```

## Provider Configuration Fields

Each provider in the `providers` section supports these common fields:

- `email`: Email address for authentication (Cloudflare, etc.)
- `password`: API token or password
- `password_file`: Path to file containing password/token
- `login_token`: Login token for token-based auth (DNSPod, etc.)
- `login_token_file`: Path to file containing login token
- `app_key`: Application key (provider-specific)
- `app_secret`: Application secret (provider-specific)
- `consumer_key`: Consumer key (provider-specific)

## Domain Configuration

Domains support an optional `provider` field:

```json
{
  "domain_name": "example.com",
  "sub_domains": ["www", "api", "@"],
  "provider": "cloudflare"
}
```

If `provider` is omitted, the domain uses the global `provider` setting.

## Supported Providers

All existing providers are supported in multi-provider mode:

- Cloudflare
- DNSPod  
- DigitalOcean
- Alidns
- Google
- HE (Hurricane Electric)
- Dreamhost
- Duck DNS
- NoIP
- Scaleway
- DynV6
- Linode
- Strato
- Loopiase
- Infomaniak
- Hetzner
- OVH
- Dynu
- IONOS
- TransIP

## Migration Guide

### From Single to Multi-Provider

1. **Keep existing config working**: No changes needed for current setups
2. **Add new providers gradually**: 
   ```json
   {
     "provider": "DNSPod",        // Keep existing
     "login_token": "old-token",
     "providers": {               // Add new providers
       "cloudflare": {
         "email": "user@example.com",
         "password": "cf-token"
       }
     },
     "domains": [
       {
         "domain_name": "old-domain.com",
         "sub_domains": ["www"]    // Uses DNSPod (global provider)
       },
       {
         "domain_name": "new-domain.com", 
         "sub_domains": ["www"],
         "provider": "cloudflare"  // Uses Cloudflare
       }
     ]
   }
   ```

3. **Complete migration**: Remove global provider once all domains specify providers

## Configuration Examples

See the example configuration files:

- `configs/config_multi_sample.json` - Full multi-provider setup
- `configs/config_multi_sample.yaml` - YAML version of multi-provider setup  
- `configs/config_legacy_compatible.json` - Shows backward compatibility
- `configs/config_mixed_sample.json` - Mixed legacy + new provider configuration

## Logging and Notifications

Multi-provider configurations include provider information in log messages and notifications:

```
INFO [2024-01-01T12:00:00Z] Initialized provider: cloudflare
INFO [2024-01-01T12:00:00Z] Initialized provider: dnspod  
INFO [2024-01-01T12:00:00Z] [ www, api ] of example.com (via cloudflare)
```

## Benefits

1. **Consolidation**: Manage multiple DNS providers from one GoDNS instance
2. **Flexibility**: Use the best provider for each domain 
3. **Redundancy**: Distribute domains across providers for resilience
4. **Migration**: Gradually move domains between providers
5. **Cost Optimization**: Use different providers based on pricing/features
6. **Compliance**: Meet requirements for geographic distribution

## Troubleshooting

### Provider Not Found Error
```
ERROR provider 'cloudflare' not found for domain example.com
```
**Solution**: Ensure the provider is configured in the `providers` section.

### Authentication Failures
```  
ERROR failed to create provider cloudflare: authentication failed
```
**Solution**: Verify credentials in the provider configuration section.

### Mixed Configuration Issues
If a domain doesn't specify a `provider` field, it will use the global `provider`. Ensure:
1. Global `provider` is set when using mixed configuration
2. All required credentials are provided for the global provider