module MarketplaceAnalytics

export forecast_revenue

"""Simple baseline forecast until production model is added."""
function forecast_revenue(monthly::Vector{Float64})
    isempty(monthly) && return 0.0
    return sum(monthly) / length(monthly)
end

end
