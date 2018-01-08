package com.kerr.domain;

import java.io.Serializable;
import org.springframework.cassandra.core.PrimaryKeyType;
import org.springframework.data.cassandra.mapping.PrimaryKeyClass;
import org.springframework.data.cassandra.mapping.PrimaryKeyColumn;

@PrimaryKeyClass
public class CampgroundKey implements Serializable{

    @PrimaryKeyColumn(name = "park_id", ordinal = 0, type = PrimaryKeyType.PARTITIONED)
    private String parkId;

    @PrimaryKeyColumn(name = "campground_id", ordinal = 1, type = PrimaryKeyType.CLUSTERED)
    private String campgroundId = "";

    public String getCampgroundId() {
        return campgroundId;
    }

    public String getParkId() {
        return parkId;
    }

    public CampgroundKey() {}

    public CampgroundKey(String parkId) {
        this(parkId, "");
    }

    public CampgroundKey(String parkId, String campgroundId) {
        this.parkId = parkId;
        this.campgroundId = campgroundId;
    }
}
