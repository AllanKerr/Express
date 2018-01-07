package com.kerr;

import java.util.Date;
import javax.validation.constraints.Min;
import javax.validation.constraints.NotNull;
import org.hibernate.validator.constraints.NotEmpty;

public final class SearchRequest {

    @NotEmpty(message = "Park is mandatory.")
    private String park;

    @NotEmpty(message = "Campground is mandatory.")
    private String campground;

    @NotNull(message = "Arrival start is mandatory.")
    private Date arrivalStart;

    @NotNull(message = "Arrival end is mandatory.")
    private Date arrivalEnd;

    @Min(1)
    private int nights;

    public String getPark() {
        return park;
    }

    public String getCampground() {
        return campground;
    }

    public Date getArrivalStart() {
        return arrivalStart;
    }

    public Date getArrivalEnd() {
        return arrivalEnd;
    }

    public int getNights() {
        return nights;
    }
}
