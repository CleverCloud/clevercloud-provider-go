$ clever addon providers
$ clever addon providers show X


POST https://api.clever-cloud.com/v2/organisations/orga_f12cf9a3-09cc-4e4d-bfa4-8ac6f915dd7b/addonproviders

{"id":"matomo-addon","name":"Matomo","website":null,"supportEmail":null,"googlePlusName":null,"twitterName":null,"analyticsId":null,"shortDesc":null,"longDesc":null,"logoUrl":null,"status":"ALPHA","openInNewTab":false,"canUpgrade":false,"regions":[],"plans":[],"features":[]}


PUT https://api.clever-cloud.com/v2/organisations/orga_f12cf9a3-09cc-4e4d-bfa4-8ac6f915dd7b/addonproviders/matomo-addon

{"id":"matomo-addon","name":"Matomo","website":"test","supportEmail":"test@test","googlePlusName":"test","twitterName":"testt","analyticsId":"test","shortDesc":"test","longDesc":"test","logoUrl":"test","status":"ALPHA","openInNewTab":false,"canUpgrade":false,"regions":[]}
