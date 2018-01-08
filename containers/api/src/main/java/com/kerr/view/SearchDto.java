package com.kerr.view;

import com.kerr.domain.Search;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import javax.validation.constraints.Min;
import javax.validation.constraints.NotNull;
import org.hibernate.validator.constraints.NotEmpty;

public final class SearchDto {

    @NotEmpty
    private List<CampgroundIdDto> campgrounds;

    @NotNull
    private Date rangeStart;

    @NotNull
    private Date rangeEnd;

    @Min(1)
    private int nights;

    public List<CampgroundIdDto> getCampgrounds() {
        return campgrounds;
    }

    public Date getRangeStart() {
        return rangeStart;
    }

    public Date getRangeEnd() {
        return rangeEnd;
    }

    public int getNights() {
        return nights;
    }

    public Iterable<Search> getSearches(String userId) {
        assert userId != null;
        List<Search> searches = new ArrayList<>(campgrounds.size());
        for (CampgroundIdDto id : campgrounds) {
            Search search = new Search(userId, id.getParkId(), id.getCampgroundId(), id.getSectionId(), rangeStart, rangeEnd, nights);
            searches.add(search);
        }
        return searches;
    }
}
