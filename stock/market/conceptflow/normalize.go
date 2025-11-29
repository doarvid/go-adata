package conceptflow

func NormalizeConceptFlows(in []ConceptFlow) []ConceptFlow {
    for i := range in {
        if in[i].MainNetInflow < 0 { in[i].MainNetInflow = 0 }
        if in[i].MaxNetInflow < 0 { in[i].MaxNetInflow = 0 }
        if in[i].LgNetInflow < 0 { in[i].LgNetInflow = 0 }
        if in[i].MidNetInflow < 0 { in[i].MidNetInflow = 0 }
        if in[i].SmNetInflow < 0 { in[i].SmNetInflow = 0 }
        if in[i].MainNetInflowRate < 0 { in[i].MainNetInflowRate = 0 }
        if in[i].MaxNetInflowRate < 0 { in[i].MaxNetInflowRate = 0 }
        if in[i].LgNetInflowRate < 0 { in[i].LgNetInflowRate = 0 }
        if in[i].MidNetInflowRate < 0 { in[i].MidNetInflowRate = 0 }
        if in[i].SmNetInflowRate < 0 { in[i].SmNetInflowRate = 0 }
    }
    return in
}

