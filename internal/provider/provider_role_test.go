package provider

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestResourcesIncludeRole(t *testing.T) {
    pfn := New("test")
    p := pfn()

    resources := p.Resources(context.Background())
    found := false
    for _, rf := range resources {
        r := rf()
        var mr resource.MetadataResponse
        r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "sda"}, &mr)
        if mr.TypeName == "sda_role" {
            found = true
            break
        }
    }

    if !found {
        t.Fatalf("sda_role resource not registered in provider resources")
    }
}
