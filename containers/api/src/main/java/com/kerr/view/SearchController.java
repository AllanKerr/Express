package com.kerr.view;

import com.kerr.domain.Search;
import com.kerr.repositories.SearchRepository;
import java.util.List;
import javax.servlet.http.HttpServletRequest;
import javax.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping(value = "/searches/v1")
public class SearchController {

    private SearchRepository searches;

    @Autowired
    public SearchController(SearchRepository searches) {
        assert searches != null;
        this.searches = searches;
    }

    @RequestMapping(value = "/add", method = RequestMethod.POST)
    public HttpStatus addSearch(@Valid @RequestBody SearchDto search, HttpServletRequest request) {

        String userId = request.getHeader("User-Id");
        if (userId == null || userId.isEmpty()) {
            return HttpStatus.UNAUTHORIZED;
        }
        searches.save(search.getSearches(userId));
        return HttpStatus.OK;
    }

    @RequestMapping(value = "/list", method = RequestMethod.GET)
    public ResponseEntity<List<Search>> listSearches(HttpServletRequest request) {

        ResponseEntity<List<Search>> response;
        String userId = request.getHeader("User-Id");
        if (userId == null || userId.isEmpty()) {
            response = ResponseEntity.status(HttpStatus.UNAUTHORIZED).build();
        } else {
            List<Search> body = searches.findAllByKey_UserId(userId);
            return new ResponseEntity<>(body, HttpStatus.OK);

        }
        return response;
    }
}
