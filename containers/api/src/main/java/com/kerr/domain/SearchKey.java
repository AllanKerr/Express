package com.kerr.domain;

import java.io.Serializable;
import org.springframework.cassandra.core.PrimaryKeyType;
import org.springframework.data.cassandra.mapping.PrimaryKeyClass;
import org.springframework.data.cassandra.mapping.PrimaryKeyColumn;

@PrimaryKeyClass
public class SearchKey implements Serializable {

    @PrimaryKeyColumn(name = "user_id", ordinal = 0, type = PrimaryKeyType.PARTITIONED)
    private String userId;

    @PrimaryKeyColumn(name = "park_id", ordinal = 1, type = PrimaryKeyType.CLUSTERED)
    private String parkId;

    @PrimaryKeyColumn(name = "campground_id", ordinal = 2, type = PrimaryKeyType.CLUSTERED)
    private String campgroundId = "";

    @PrimaryKeyColumn(name = "section_id", ordinal = 3, type = PrimaryKeyType.CLUSTERED)
    private String sectionId = "";

    public String getCampgroundId() {
        return campgroundId;
    }

    public String getParkId() {
        return parkId;
    }

    public String getSectionId() {
        return sectionId;
    }

    public String getUserId() {
        return userId;
    }

    public SearchKey() {}

    public SearchKey(String userId, String parkId) {
        this(userId, parkId, "");
    }

    public SearchKey(String userId, String parkId, String campgroundId) {
        this(userId, parkId, campgroundId, "");
    }

    public SearchKey(String userId, String parkId, String campgroundId, String sectionId) {
        assert parkId != null;
        if (campgroundId == null) {
            campgroundId = "";
        }
        if (sectionId == null) {
            sectionId = "";
        }
        this.parkId = parkId;
        this.campgroundId = campgroundId;
        this.sectionId = sectionId;
        this.userId = userId;
    }
}
