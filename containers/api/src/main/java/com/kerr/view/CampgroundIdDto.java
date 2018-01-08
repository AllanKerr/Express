package com.kerr.view;

import org.hibernate.validator.constraints.NotEmpty;

public class CampgroundIdDto {

    @NotEmpty
    private String parkId;

    private String campgroundId;

    private String sectionId;

    public String getParkId() {
        return parkId;
    }

    public String getCampgroundId() {
        return campgroundId;
    }

    public String getSectionId() {
        return sectionId;
    }
}
