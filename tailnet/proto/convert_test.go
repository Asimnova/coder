package proto_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	"tailscale.com/tailcfg"

	"github.com/coder/coder/v2/tailnet/proto"
)

// tailscaleDERPMap on 2023-11-20 for testing purposes.
const tailscaleDERPMap = `{"Regions":{"1":{"RegionID":1,"RegionCode":"nyc","RegionName":"New York City","Nodes":[{"Name":"1f","RegionID":1,"HostName":"derp1f.tailscale.com","IPv4":"199.38.181.104","IPv6":"2607:f740:f::bc"},{"Name":"1g","RegionID":1,"HostName":"derp1g.tailscale.com","IPv4":"209.177.145.120","IPv6":"2607:f740:f::3eb"},{"Name":"1h","RegionID":1,"HostName":"derp1h.tailscale.com","IPv4":"199.38.181.93","IPv6":"2607:f740:f::afd"},{"Name":"1i","RegionID":1,"HostName":"derp1i.tailscale.com","IPv4":"199.38.181.103","IPv6":"2607:f740:f::e19","STUNOnly":true}]},"10":{"RegionID":10,"RegionCode":"sea","RegionName":"Seattle","Nodes":[{"Name":"10b","RegionID":10,"HostName":"derp10b.tailscale.com","IPv4":"192.73.240.161","IPv6":"2607:f740:14::61c"},{"Name":"10c","RegionID":10,"HostName":"derp10c.tailscale.com","IPv4":"192.73.240.121","IPv6":"2607:f740:14::40c"},{"Name":"10d","RegionID":10,"HostName":"derp10d.tailscale.com","IPv4":"192.73.240.132","IPv6":"2607:f740:14::500"}]},"11":{"RegionID":11,"RegionCode":"sao","RegionName":"São Paulo","Nodes":[{"Name":"11b","RegionID":11,"HostName":"derp11b.tailscale.com","IPv4":"148.163.220.129","IPv6":"2607:f740:1::211"},{"Name":"11c","RegionID":11,"HostName":"derp11c.tailscale.com","IPv4":"148.163.220.134","IPv6":"2607:f740:1::861"},{"Name":"11d","RegionID":11,"HostName":"derp11d.tailscale.com","IPv4":"148.163.220.210","IPv6":"2607:f740:1::2e6"}]},"12":{"RegionID":12,"RegionCode":"ord","RegionName":"Chicago","Nodes":[{"Name":"12d","RegionID":12,"HostName":"derp12d.tailscale.com","IPv4":"209.177.158.246","IPv6":"2607:f740:e::811"},{"Name":"12e","RegionID":12,"HostName":"derp12e.tailscale.com","IPv4":"209.177.158.15","IPv6":"2607:f740:e::b17"},{"Name":"12f","RegionID":12,"HostName":"derp12f.tailscale.com","IPv4":"199.38.182.118","IPv6":"2607:f740:e::4c8"}]},"13":{"RegionID":13,"RegionCode":"den","RegionName":"Denver","Nodes":[{"Name":"13b","RegionID":13,"HostName":"derp13b.tailscale.com","IPv4":"192.73.242.187","IPv6":"2607:f740:16::640"},{"Name":"13c","RegionID":13,"HostName":"derp13c.tailscale.com","IPv4":"192.73.242.28","IPv6":"2607:f740:16::5c"},{"Name":"13d","RegionID":13,"HostName":"derp13d.tailscale.com","IPv4":"192.73.242.204","IPv6":"2607:f740:16::c23"}]},"14":{"RegionID":14,"RegionCode":"ams","RegionName":"Amsterdam","Nodes":[{"Name":"14b","RegionID":14,"HostName":"derp14b.tailscale.com","IPv4":"176.58.93.248","IPv6":"2a00:dd80:3c::807"},{"Name":"14c","RegionID":14,"HostName":"derp14c.tailscale.com","IPv4":"176.58.93.147","IPv6":"2a00:dd80:3c::b09"},{"Name":"14d","RegionID":14,"HostName":"derp14d.tailscale.com","IPv4":"176.58.93.154","IPv6":"2a00:dd80:3c::3d5"}]},"15":{"RegionID":15,"RegionCode":"jnb","RegionName":"Johannesburg","Nodes":[{"Name":"15b","RegionID":15,"HostName":"derp15b.tailscale.com","IPv4":"102.67.165.90","IPv6":"2c0f:edb0:0:10::963"},{"Name":"15c","RegionID":15,"HostName":"derp15c.tailscale.com","IPv4":"102.67.165.185","IPv6":"2c0f:edb0:0:10::b59"},{"Name":"15d","RegionID":15,"HostName":"derp15d.tailscale.com","IPv4":"102.67.165.36","IPv6":"2c0f:edb0:0:10::599"}]},"16":{"RegionID":16,"RegionCode":"mia","RegionName":"Miami","Nodes":[{"Name":"16b","RegionID":16,"HostName":"derp16b.tailscale.com","IPv4":"192.73.243.135","IPv6":"2607:f740:17::476"},{"Name":"16c","RegionID":16,"HostName":"derp16c.tailscale.com","IPv4":"192.73.243.229","IPv6":"2607:f740:17::4e4"},{"Name":"16d","RegionID":16,"HostName":"derp16d.tailscale.com","IPv4":"192.73.243.141","IPv6":"2607:f740:17::475"}]},"17":{"RegionID":17,"RegionCode":"lax","RegionName":"Los Angeles","Nodes":[{"Name":"17b","RegionID":17,"HostName":"derp17b.tailscale.com","IPv4":"192.73.244.245","IPv6":"2607:f740:c::646"},{"Name":"17c","RegionID":17,"HostName":"derp17c.tailscale.com","IPv4":"208.111.40.12","IPv6":"2607:f740:c::10"},{"Name":"17d","RegionID":17,"HostName":"derp17d.tailscale.com","IPv4":"208.111.40.216","IPv6":"2607:f740:c::e1b"}]},"18":{"RegionID":18,"RegionCode":"par","RegionName":"Paris","Nodes":[{"Name":"18b","RegionID":18,"HostName":"derp18b.tailscale.com","IPv4":"176.58.90.147","IPv6":"2a00:dd80:3e::363"},{"Name":"18c","RegionID":18,"HostName":"derp18c.tailscale.com","IPv4":"176.58.90.207","IPv6":"2a00:dd80:3e::c19"},{"Name":"18d","RegionID":18,"HostName":"derp18d.tailscale.com","IPv4":"176.58.90.104","IPv6":"2a00:dd80:3e::f2e"}]},"19":{"RegionID":19,"RegionCode":"mad","RegionName":"Madrid","Nodes":[{"Name":"19b","RegionID":19,"HostName":"derp19b.tailscale.com","IPv4":"45.159.97.144","IPv6":"2a00:dd80:14:10::335"},{"Name":"19c","RegionID":19,"HostName":"derp19c.tailscale.com","IPv4":"45.159.97.61","IPv6":"2a00:dd80:14:10::20"},{"Name":"19d","RegionID":19,"HostName":"derp19d.tailscale.com","IPv4":"45.159.97.233","IPv6":"2a00:dd80:14:10::34a"}]},"2":{"RegionID":2,"RegionCode":"sfo","RegionName":"San Francisco","Nodes":[{"Name":"2d","RegionID":2,"HostName":"derp2d.tailscale.com","IPv4":"192.73.252.65","IPv6":"2607:f740:0:3f::287"},{"Name":"2e","RegionID":2,"HostName":"derp2e.tailscale.com","IPv4":"192.73.252.134","IPv6":"2607:f740:0:3f::44c"},{"Name":"2f","RegionID":2,"HostName":"derp2f.tailscale.com","IPv4":"208.111.34.178","IPv6":"2607:f740:0:3f::f4"}]},"20":{"RegionID":20,"RegionCode":"hkg","RegionName":"Hong Kong","Nodes":[{"Name":"20b","RegionID":20,"HostName":"derp20b.tailscale.com","IPv4":"103.6.84.152","IPv6":"2403:2500:8000:1::ef6"},{"Name":"20c","RegionID":20,"HostName":"derp20c.tailscale.com","IPv4":"205.147.105.30","IPv6":"2403:2500:8000:1::5fb"},{"Name":"20d","RegionID":20,"HostName":"derp20d.tailscale.com","IPv4":"205.147.105.78","IPv6":"2403:2500:8000:1::e9a"}]},"21":{"RegionID":21,"RegionCode":"tor","RegionName":"Toronto","Nodes":[{"Name":"21b","RegionID":21,"HostName":"derp21b.tailscale.com","IPv4":"162.248.221.199","IPv6":"2607:f740:50::1d1"},{"Name":"21c","RegionID":21,"HostName":"derp21c.tailscale.com","IPv4":"162.248.221.215","IPv6":"2607:f740:50::f10"},{"Name":"21d","RegionID":21,"HostName":"derp21d.tailscale.com","IPv4":"162.248.221.248","IPv6":"2607:f740:50::ca4"}]},"22":{"RegionID":22,"RegionCode":"waw","RegionName":"Warsaw","Nodes":[{"Name":"22b","RegionID":22,"HostName":"derp22b.tailscale.com","IPv4":"45.159.98.196","IPv6":"2a00:dd80:40:100::316"},{"Name":"22c","RegionID":22,"HostName":"derp22c.tailscale.com","IPv4":"45.159.98.253","IPv6":"2a00:dd80:40:100::3f"},{"Name":"22d","RegionID":22,"HostName":"derp22d.tailscale.com","IPv4":"45.159.98.145","IPv6":"2a00:dd80:40:100::211"}]},"23":{"RegionID":23,"RegionCode":"dbi","RegionName":"Dubai","Nodes":[{"Name":"23b","RegionID":23,"HostName":"derp23b.tailscale.com","IPv4":"185.34.3.232","IPv6":"2a00:dd80:3f:100::76f"},{"Name":"23c","RegionID":23,"HostName":"derp23c.tailscale.com","IPv4":"185.34.3.207","IPv6":"2a00:dd80:3f:100::a50"},{"Name":"23d","RegionID":23,"HostName":"derp23d.tailscale.com","IPv4":"185.34.3.75","IPv6":"2a00:dd80:3f:100::97e"}]},"24":{"RegionID":24,"RegionCode":"hnl","RegionName":"Honolulu","Nodes":[{"Name":"24b","RegionID":24,"HostName":"derp24b.tailscale.com","IPv4":"208.83.234.151","IPv6":"2001:19f0:c000:c586:5400:04ff:fe26:2ba6"},{"Name":"24c","RegionID":24,"HostName":"derp24c.tailscale.com","IPv4":"208.83.233.233","IPv6":"2001:19f0:c000:c591:5400:04ff:fe26:2c5f"},{"Name":"24d","RegionID":24,"HostName":"derp24d.tailscale.com","IPv4":"208.72.155.133","IPv6":"2001:19f0:c000:c564:5400:04ff:fe26:2ba8"}]},"25":{"RegionID":25,"RegionCode":"nai","RegionName":"Nairobi","Nodes":[{"Name":"25b","RegionID":25,"HostName":"derp25b.tailscale.com","IPv4":"102.67.167.245","IPv6":"2c0f:edb0:2000:1::2e9"},{"Name":"25c","RegionID":25,"HostName":"derp25c.tailscale.com","IPv4":"102.67.167.37","IPv6":"2c0f:edb0:2000:1::2c7"},{"Name":"25d","RegionID":25,"HostName":"derp25d.tailscale.com","IPv4":"102.67.167.188","IPv6":"2c0f:edb0:2000:1::188"}]},"3":{"RegionID":3,"RegionCode":"sin","RegionName":"Singapore","Nodes":[{"Name":"3b","RegionID":3,"HostName":"derp3b.tailscale.com","IPv4":"43.245.49.105","IPv6":"2403:2500:300::b0c"},{"Name":"3c","RegionID":3,"HostName":"derp3c.tailscale.com","IPv4":"43.245.49.83","IPv6":"2403:2500:300::57a"},{"Name":"3d","RegionID":3,"HostName":"derp3d.tailscale.com","IPv4":"43.245.49.144","IPv6":"2403:2500:300::df9"}]},"4":{"RegionID":4,"RegionCode":"fra","RegionName":"Frankfurt","Nodes":[{"Name":"4f","RegionID":4,"HostName":"derp4f.tailscale.com","IPv4":"185.40.234.219","IPv6":"2a00:dd80:20::a25"},{"Name":"4g","RegionID":4,"HostName":"derp4g.tailscale.com","IPv4":"185.40.234.113","IPv6":"2a00:dd80:20::8f"},{"Name":"4h","RegionID":4,"HostName":"derp4h.tailscale.com","IPv4":"185.40.234.77","IPv6":"2a00:dd80:20::bcf"}]},"5":{"RegionID":5,"RegionCode":"syd","RegionName":"Sydney","Nodes":[{"Name":"5b","RegionID":5,"HostName":"derp5b.tailscale.com","IPv4":"43.245.48.220","IPv6":"2403:2500:9000:1::ce7"},{"Name":"5c","RegionID":5,"HostName":"derp5c.tailscale.com","IPv4":"43.245.48.50","IPv6":"2403:2500:9000:1::f57"},{"Name":"5d","RegionID":5,"HostName":"derp5d.tailscale.com","IPv4":"43.245.48.250","IPv6":"2403:2500:9000:1::43"}]},"6":{"RegionID":6,"RegionCode":"blr","RegionName":"Bangalore","Nodes":[{"Name":"6a","RegionID":6,"HostName":"derp6.tailscale.com","IPv4":"68.183.90.120","IPv6":"2400:6180:100:d0::982:d001"}]},"7":{"RegionID":7,"RegionCode":"tok","RegionName":"Tokyo","Nodes":[{"Name":"7b","RegionID":7,"HostName":"derp7b.tailscale.com","IPv4":"103.84.155.178","IPv6":"2403:2500:400:20::b79"},{"Name":"7c","RegionID":7,"HostName":"derp7c.tailscale.com","IPv4":"103.84.155.188","IPv6":"2403:2500:400:20::835"},{"Name":"7d","RegionID":7,"HostName":"derp7d.tailscale.com","IPv4":"103.84.155.46","IPv6":"2403:2500:400:20::cfe"}]},"8":{"RegionID":8,"RegionCode":"lhr","RegionName":"London","Nodes":[{"Name":"8e","RegionID":8,"HostName":"derp8e.tailscale.com","IPv4":"176.58.92.144","IPv6":"2a00:dd80:3a::b33"},{"Name":"8f","RegionID":8,"HostName":"derp8f.tailscale.com","IPv4":"176.58.88.183","IPv6":"2a00:dd80:3a::dfa"},{"Name":"8g","RegionID":8,"HostName":"derp8g.tailscale.com","IPv4":"176.58.92.254","IPv6":"2a00:dd80:3a::ed"}]},"9":{"RegionID":9,"RegionCode":"dfw","RegionName":"Dallas","Nodes":[{"Name":"9d","RegionID":9,"HostName":"derp9d.tailscale.com","IPv4":"209.177.156.94","IPv6":"2607:f740:100::c05"},{"Name":"9e","RegionID":9,"HostName":"derp9e.tailscale.com","IPv4":"192.73.248.83","IPv6":"2607:f740:100::359"},{"Name":"9f","RegionID":9,"HostName":"derp9f.tailscale.com","IPv4":"209.177.156.197","IPv6":"2607:f740:100::cad"}]}}}`

func TestDERPMap(t *testing.T) {
	t.Parallel()

	derpMap := &tailcfg.DERPMap{}
	err := json.Unmarshal([]byte(tailscaleDERPMap), derpMap)
	require.NoError(t, err)

	// The tailscale DERPMap doesn't have HomeParams.
	derpMap.HomeParams = &tailcfg.DERPHomeParams{
		RegionScore: map[int]float64{
			1: 2,
			2: 3,
		},
	}

	// Add a region and node that uses every single field.
	derpMap.Regions[999] = &tailcfg.DERPRegion{
		RegionID:      999,
		EmbeddedRelay: true,
		RegionCode:    "zzz",
		RegionName:    "Cool Region",
		Avoid:         true,

		Nodes: []*tailcfg.DERPNode{
			{
				Name:             "zzz1",
				RegionID:         999,
				HostName:         "coolderp.com",
				IPv4:             "1.2.3.4",
				IPv6:             "2001:db8::1",
				STUNPort:         1234,
				STUNOnly:         true,
				DERPPort:         5678,
				InsecureForTests: true,
				ForceHTTP:        true,
				STUNTestIP:       "5.6.7.8",
				CanPort80:        true,
			},
		},
	}

	protoMap := proto.DERPMapToProto(derpMap)
	require.NotNil(t, protoMap)

	derpMap2 := proto.DERPMapFromProto(protoMap)
	require.NotNil(t, derpMap2)

	require.Equal(t, derpMap, derpMap2)
}
