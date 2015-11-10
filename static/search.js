$(function () {
    var queryString = $('#queryString');
    var results = $('#results');

    $('#queryString').focus();
    $('#myTab li:eq(0) a').tab('show');
    $('#showResults').hide("fast");

    $("#searchBtn").click(function() {
        searchPost();
    });

    $("#clearResults").click(function () {
        var results = $('#results');
        results.html("");
    });

    /* START: The following taken from https://github.com/Reactive-Extensions/RxJS */
    /* Only get the value from each key up */
    var keyups = Rx.Observable.fromEvent(queryString, 'keyup')
        .map(function (e) {
            return e.target.value;
        })
        .filter(function (text) {
            return text.length > 2;
        })
        .catch(function (e) {
            alert("error: " + e);
            return Rx.Observable.empty()
        });

    /* Now debounce the input for 500ms */
    var debounced = keyups
        .debounce(500 /* ms */);

    /* Now get only distinct values, so we eliminate the arrows and other control characters */
    var distinct = debounced
        .distinctUntilChanged();
    /* END https://github.com/Reactive-Extensions/RxJS */

    function searchAutocomplete(term) {
        return $.ajax({
            url: '/search/' + $('#queryString').val(),
            dataType: 'json',
            data: {
                format: 'json',
                queryString: term
            }
        }).promise();
    }

    var searchPost = function () {
        $.ajax({
            url: '/searchPost/' + $('#queryString').val()
            , type: 'POST'
            , dataType: 'json'
            , processData: false
            , success: function (json, statusText, xhr) {
                successSearchResults(data);
                showFacets(data);
            }
            , error: function (xhr, message, error) {
                console.error("Error while loading data from ElasticSearch", message);
                alert(error);
            }
        });
    };

    var suggestions = distinct
        .flatMapLatest(searchAutocomplete);

    suggestions.forEach(
        function (data) {
            $('#showResults').show("fast");
            var tookDiv = $('#tookDiv');

            var took = data.results.took;
            if (data.results != -1) {
                tookDiv.html("");
                tookDiv.append("Found  <b>" + data.results.hits.total + " </b>documents in " + took + "ms" + " for query " + "<b>" + $('#queryString').val() + "</b>");
                tookDiv.append("<br /><br />");

                successSearchResults(data);
                showFacets(data);

                /*
                results
                    .empty()
                    .append($.map(data.results, function (value) {
                        return $('<li>').text(value.author);
                    }));
                    */
            }
        },
        function (error) {
            results
                .empty()
                .append($('<li>'))
                .text('Error:' + error);
        });

});