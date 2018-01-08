package com.kerr.view;

import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/searches/v1")
public class SearchController {


    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public HttpStatus addSearch(@Valid @RequestBody SearchDto search, HttpServletRequest request) {

        String userId = request.getHeader("User-Id");

        System.out.println("Nights: " + search.getNights() + "\nStart:" + search.getRangeStart() + "\nEnd:" + search.getRangeEnd() + "\n");
        for (CampgroundIdDto id : search.getCampgrounds()) {
            System.out.println("Id: " + id.getParkId() + " camp: " + id.getCampgroundId() + " sec: " + id.getSectionId() + "\n");
        }


        return HttpStatus.OK;
    }
}
