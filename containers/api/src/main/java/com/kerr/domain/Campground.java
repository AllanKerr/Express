package com.kerr.domain;

import org.springframework.data.cassandra.mapping.Column;
import org.springframework.data.cassandra.mapping.PrimaryKey;
import org.springframework.data.cassandra.mapping.Table;

@Table("campgrounds")
public class Campground {

    @PrimaryKey
    private CampgroundKey key;

    private String organization;

    @Column("park_name")
    private String parkName;

    @Column("campground_name")
    private String campgroundName;

    public String getParkId() {
        return key.getParkId();
    }

    public String getCampgroundId() {
        return key.getCampgroundId();
    }

    public String getOrganization() {
        return organization;
    }

    public String getCampgroundName() {
        return campgroundName;
    }

    public void setCampgroundName(String campgroundName) {
        this.campgroundName = campgroundName;
    }

    public String getParkName() {
        return parkName;
    }

    public void setParkName(String parkName) {
        this.parkName = parkName;
    }

    public Campground(){}

    public Campground(String parkId, String campgroundId, String organization, String parkName,
            String campgroundName) {
        this.key = new CampgroundKey(parkId, campgroundId);
        this.organization = organization;
        this.parkName = parkName;
        this.campgroundName = campgroundName;
    }
}
