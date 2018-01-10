package com.kerr.domain;

import java.util.Date;
import org.springframework.data.cassandra.mapping.Column;
import org.springframework.data.cassandra.mapping.PrimaryKey;
import org.springframework.data.cassandra.mapping.Table;

@Table("searches")
public class Search {

    @PrimaryKey
    private SearchKey key;

    @Column("range_start")
    private Date rangeStart;

    @Column("range_end")
    private Date rangeEnd;

    private int nights;

    public String getUserId() {
        return key.getUserId();
    }

    public String getParkId() {
        return key.getParkId();
    }

    public String getCampgroundId() {
        return key.getCampgroundId();
    }

    public String getSectionId() {
        return key.getSectionId();
    }

    public Date getRangeStart() {
        return rangeStart;
    }

    public void setRangeStart(Date rangeStart) {
        this.rangeStart = rangeStart;
    }

    public Date getRangeEnd() {
        return rangeEnd;
    }

    public void setRangeEnd(Date rangeEnd) {
        this.rangeEnd = rangeEnd;
    }

    public int getNights() {
        return nights;
    }

    public void setNights(int nights) {
        this.nights = nights;
    }

    public Search() {}

    public Search(SearchKey key, Date rangeStart, Date rangeEnd, int nights) {
        this.key = key;
        this.rangeStart = rangeStart;
        this.rangeEnd = rangeEnd;
        this.nights = nights;
    }

    public Search(String userId, String parkId, String campgroundId, String sectionId, Date rangeStart, Date rangeEnd, int nights) {
        this.key = new SearchKey(userId, parkId, campgroundId, sectionId);
        this.rangeStart = rangeStart;
        this.rangeEnd = rangeEnd;
        this.nights = nights;
    }
}
