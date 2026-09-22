output "location" {
  description = "Map of location information. Keys match the input map keys, values contain name, display_name, short_name, regional_display_name, and paired_region_name."
  value       = local.location
}
