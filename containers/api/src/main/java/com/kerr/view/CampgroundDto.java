package com.kerr.view;

import com.kerr.domain.Campground;
import org.hibernate.validator.constraints.NotEmpty;

public final class CampgroundDto {

    public String getParkId() {
        return parkId;
    }

    public String getCampgroundId() {
        return campgroundId;
    }

    public String getOrganization() {
        return organization;
    }

    public String getParkName() {
        return parkName;
    }

    public String getCampgroundName() {
        return campgroundName;
    }

    @NotEmpty

    private String parkId;

    @NotEmpty
    private String campgroundId;

    @NotEmpty
    private String organization;

    @NotEmpty
    private String parkName;

    private String campgroundName;

    public Campground getCampground() {
        return new Campground(parkId, campgroundId, organization, parkName, parkName);
    }
}
