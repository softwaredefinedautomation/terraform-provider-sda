package provider

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestResourcesIncludeUser(t *testing.T) {
    pfn := New("test")
    p := pfn()

    resources := p.Resources(context.Background())
    found := false
    for _, rf := range resources {
        r := rf()
        // call Metadata to get TypeName
        var mr resource.MetadataResponse
        r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "sda"}, &mr)
        if mr.TypeName == "sda_user" {
            found = true
            break
        }
    }

    if !found {
        t.Fatalf("sda_user resource not registered in provider resources")
    }
}
