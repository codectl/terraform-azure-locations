module "regions" {
  source  = "cloudnationhq/locations/azure"
  version = "~> 1.0"

  location = {
    primary   = "westeurope"
    secondary = "North Europe"
    backup    = "eastus"
  }
}

output "all_locations" {
  description = "Complete information for all locations"
  value       = module.regions.location
}

output "primary_name" {
  description = "Standard Azure region name for primary location"
  value       = module.regions.location.primary.name
}

output "primary_display_name" {
  description = "Human-readable region name for primary location"
  value       = module.regions.location.primary.display_name
}

output "secondary_short_name" {
  description = "Abbreviated region name for secondary location"
  value       = module.regions.location.secondary.short_name
}

output "backup_paired_region" {
  description = "Azure paired region name for backup location"
  value       = module.regions.location.backup.paired_region_name
}
