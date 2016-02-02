var successSearchResults = function (json) {
    if (json.results.hits.total > 0) {
        var items = json.results.hits.hits;
        var observationTableDiv = $('#observationTableDiv');

        $("#observationTableDiv > tbody").empty();
        $.each(items, function () {
            var metaData = "";
            var highlights = this.highlight;

            $.each(this._source.annotations, function (key, val) {
                var hdata = doHighlight(highlights, key);

                if (hdata == null) {
                    metaData += key + " : " + val;
                } else {
                    metaData += hdata;
                }

                metaData += "<br/>";
            });

            var row = "<tr><td>" +
                this._id + "</td><td>" +
                doSampleHighlight(highlights, "name", this._source.name) + "</td><td>" +
                doSampleHighlight(highlights, "genome", this._source.genome) + "</td><td>" +
                doSampleHighlight(highlights, "gender", this._source.gender) + "</td><td>" +
                doSampleHighlight(highlights, "variants", this._source.variants) + "</td><td>" +
                doSampleHighlight(highlights, "timeCreated", this._source.timeCreated) + "</td><td>" +
                metaData + "</td></tr>";
            $("#observationTableDiv tbody").append(row);
        });
    }
}

var showFacets = function (json) {
    if (json.results.hits.total > 0) {
        var genomeItems = json.results.facets.genome.terms;
        var genderItems = json.results.facets.gender.terms;
        var ethnicityItems = json.results.facets.ethnicity.terms;

        var topResults = json.results.hits.hits;

        var genomeListDiv = $('#genomeListDiv');
        var genderListDiv = $('#genderListDiv');
        var ethnicityListDiv = $('#ethnicityListDiv');
        var topResultsDiv = $('#topResultsDiv');

        genomeListDiv.html("");
        genderListDiv.html("");
        ethnicityListDiv.html("");
        topResultsDiv.html("");
        $.each(genomeItems, function () {
            genomeListDiv.append(
                $(document.createElement('li')).html(this.term + " (" + this.count + ")")
            );
        });
        $.each(genderItems, function () {
            genderListDiv.append(
                $(document.createElement('li')).html(this.term + " (" + this.count + ")")
            );
        });
        $.each(ethnicityItems, function () {
            ethnicityListDiv.append(
                $(document.createElement('li')).html(this.term + " (" + this.count + ")")
            );
        });
    }
};

var doHighlight = function (highlights, key) {
    var hdata = null;

    if (highlights == null) {
        return null;
    }

    $.each(highlights, function (hkey, hval) {
        if (hkey == "annotations." + key) {
            hdata = key + " : " + hval[0];
            return false;
        }
    });

    return hdata;
}

var doSampleHighlight = function (highlights, fieldName, fieldValue) {
    var hdata = fieldValue;

    if (highlights == null) {
        return hdata;
    }

    $.each(highlights, function (hkey, hval) {
        if (hkey == fieldName) {
            hdata = hval[0];
            return false;
        }
    });

    return hdata;
}
