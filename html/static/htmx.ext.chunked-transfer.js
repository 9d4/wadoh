// code source: https://github.com/JEBailey/htmx/blob/master/src/ext/chunked.js
/*
Chunked Transfer Encoding Support
======================================
This offers support for the chunked transfer encoding
*/

(function () {
  "use strict";
  var api;

  htmx.defineExtension("chunked", {
    init: function (apiRef) {
      api = apiRef;
    },

    onEvent: function (name, evt) {
      var elt = evt.target,
        xhr = evt.detail.xhr,
        last_index = 0,
        tasks = [],
        onload;

      if (xhr) {
        onload = xhr.onload;
      }

      if (name === "htmx:beforeRequest") {
        xhr.onprogress = function () {
          var chunked =
              xhr.getResponseHeader("Transfer-Encoding") === "chunked",
            current_index = xhr.responseText.length;

          if (!chunked || last_index === current_index) {
            return;
          }

          var responseText = xhr.responseText.substring(
              last_index,
              current_index
            ),
            style = api.getSwapSpecification(elt).swapStyle,
            target = api.getTarget(elt),
            settleInfo = api.makeSettleInfo(elt);

          api.withExtensions(elt, function (extension) {
            responseText = extension.transformResponse(responseText, xhr, elt);
          });

          last_index = current_index;

          api.selectAndSwap(style, target, elt, responseText, settleInfo);

          tasks = tasks.concat(settleInfo.tasks);
        };

        xhr.onload = function () {
          api.settleImmediately(tasks);
          onload.apply();
        };
      }
    },
  });
})();

// index.ts
(function() {
  let api;
  htmx.defineExtension("chunked-transfer", {
    init: function(apiRef) {
      api = apiRef;
    },
    onEvent: function(name, evt) {
      const elt = evt.target;
      if (name === "htmx:beforeRequest") {
        const xhr = evt.detail.xhr;
        xhr.onprogress = function() {
          const is_chunked = xhr.getResponseHeader("Transfer-Encoding") === "chunked";
          if (!is_chunked)
            return;
          let response = xhr.response;
          api.withExtensions(elt, function(extension) {
            if (!extension.transformResponse)
              return;
            response = extension.transformResponse(response, xhr, elt);
          });
          var swapSpec = api.getSwapSpecification(elt);
          var target = api.getTarget(elt);
          var settleInfo = api.makeSettleInfo(elt);
          api.selectAndSwap(swapSpec.swapStyle, target, elt, response, settleInfo);
          api.settleImmediately(settleInfo.tasks);
        };
      }
    }
  });
})();
