locals {
  locations_data = yamldecode(file("${path.module}/locations.yaml"))

  locations_name = {
    for location in local.locations_data.locations : location.name => {
      name                  = location.name
      display_name          = location.displayName
      short_name            = location.shortName
      regional_display_name = location.regionalDisplayName
      paired_region_name    = location.pairedRegionName
    }
  }

  locations_display_name = {
    for location in local.locations_data.locations : location.displayName => {
      name                  = location.name
      display_name          = location.displayName
      short_name            = location.shortName
      regional_display_name = location.regionalDisplayName
      paired_region_name    = location.pairedRegionName
    }
  }

  location = {
    for key, loc in var.location : key => try(coalesce(
      lookup(local.locations_name, loc, null),
      lookup(local.locations_display_name, loc, null)
      ), {
      name                  = "none"
      display_name          = "none"
      short_name            = "none"
      regional_display_name = "none"
      paired_region_name    = "none"
    })
  }
}
