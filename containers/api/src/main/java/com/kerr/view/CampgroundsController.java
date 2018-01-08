package com.kerr.view;

import com.kerr.domain.Campground;
import com.kerr.repositories.CampgroundRepository;
import javax.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/campgrounds/v1")
public class CampgroundsController {

    private CampgroundRepository campgrounds;

    @Autowired
    public CampgroundsController(CampgroundRepository parks) {
        assert parks != null;
        this.campgrounds = parks;
    }

    @RequestMapping(value = "/list", method = RequestMethod.GET)
    public Iterable<Campground> list() {
        return campgrounds.findAll();
    }

    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public HttpStatus addCampground(@Valid @RequestBody CampgroundDto campground) {
        campgrounds.insert(campground.getCampground());
        return HttpStatus.OK;
    }
}
