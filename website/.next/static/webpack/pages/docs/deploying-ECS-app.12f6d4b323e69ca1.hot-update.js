"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
self["webpackHotUpdate_N_E"]("pages/docs/deploying-ECS-app",{

/***/ "./utilities/sideBarMenu.js":
/*!**********************************!*\
  !*** ./utilities/sideBarMenu.js ***!
  \**********************************/
/***/ (function(module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"sideBarMenu\": function() { return /* binding */ sideBarMenu; }\n/* harmony export */ });\n/* module decorator */ module = __webpack_require__.hmd(module);\nvar sideBarMenu = {\n    mainMenu: [\n        {\n            title: \"getting started\",\n            nestedItems: [\n                \"welcome\",\n                \"installation\"\n            ]\n        },\n        {\n            title: \"using ize\",\n            nestedItems: [\n                \"deploying ECS app\",\n                \"deploying serverless app\"\n            ]\n        }\n    ]\n};\n\n\n;\n    var _a, _b;\n    // Legacy CSS implementations will `eval` browser code in a Node.js context\n    // to extract CSS. For backwards compatibility, we need to check we're in a\n    // browser context before continuing.\n    if (typeof self !== 'undefined' &&\n        // AMP / No-JS mode does not inject these helpers:\n        '$RefreshHelpers$' in self) {\n        var currentExports = module.__proto__.exports;\n        var prevExports = (_b = (_a = module.hot.data) === null || _a === void 0 ? void 0 : _a.prevExports) !== null && _b !== void 0 ? _b : null;\n        // This cannot happen in MainTemplate because the exports mismatch between\n        // templating and execution.\n        self.$RefreshHelpers$.registerExportsForReactRefresh(currentExports, module.id);\n        // A module can be accepted automatically based on its exports, e.g. when\n        // it is a Refresh Boundary.\n        if (self.$RefreshHelpers$.isReactRefreshBoundary(currentExports)) {\n            // Save the previous exports on update so we can compare the boundary\n            // signatures.\n            module.hot.dispose(function (data) {\n                data.prevExports = currentExports;\n            });\n            // Unconditionally accept an update to this module, we'll check if it's\n            // still a Refresh Boundary later.\n            module.hot.accept();\n            // This field is set when the previous version of this module was a\n            // Refresh Boundary, letting us know we need to check for invalidation or\n            // enqueue an update.\n            if (prevExports !== null) {\n                // A boundary can become ineligible if its exports are incompatible\n                // with the previous exports.\n                //\n                // For example, if you add/remove/change exports, we'll want to\n                // re-execute the importing modules, and force those components to\n                // re-render. Similarly, if you convert a class component to a\n                // function, we want to invalidate the boundary.\n                if (self.$RefreshHelpers$.shouldInvalidateReactRefreshBoundary(prevExports, currentExports)) {\n                    module.hot.invalidate();\n                }\n                else {\n                    self.$RefreshHelpers$.scheduleUpdate();\n                }\n            }\n        }\n        else {\n            // Since we just executed the code for the module, it's possible that the\n            // new exports made it ineligible for being a boundary.\n            // We only care about the case when we were _previously_ a boundary,\n            // because we already accepted this update (accidental side effect).\n            var isNoLongerABoundary = prevExports !== null;\n            if (isNoLongerABoundary) {\n                module.hot.invalidate();\n            }\n        }\n    }\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi91dGlsaXRpZXMvc2lkZUJhck1lbnUuanMuanMiLCJtYXBwaW5ncyI6Ijs7Ozs7QUFBTyxHQUFLLENBQUNBLFdBQVcsR0FBRyxDQUFDO0lBQ3hCQyxRQUFRLEVBQUUsQ0FBQztRQUNQLENBQUM7WUFDR0MsS0FBSyxFQUFFLENBQWlCO1lBQ3hCQyxXQUFXLEVBQUUsQ0FBQztnQkFDVixDQUFTO2dCQUNULENBQWM7WUFDbEIsQ0FBQztRQUNMLENBQUM7UUFDRCxDQUFDO1lBQ0dELEtBQUssRUFBRSxDQUFXO1lBQ2xCQyxXQUFXLEVBQUUsQ0FBQztnQkFDVixDQUFtQjtnQkFDbkIsQ0FBMEI7WUFDOUIsQ0FBQztRQUNMLENBQUM7SUFDTCxDQUFDO0FBRUwsQ0FBQyIsInNvdXJjZXMiOlsid2VicGFjazovL19OX0UvLi91dGlsaXRpZXMvc2lkZUJhck1lbnUuanM/YWZkMCJdLCJzb3VyY2VzQ29udGVudCI6WyJleHBvcnQgY29uc3Qgc2lkZUJhck1lbnUgPSB7XHJcbiAgICBtYWluTWVudTogW1xyXG4gICAgICAgIHtcclxuICAgICAgICAgICAgdGl0bGU6IFwiZ2V0dGluZyBzdGFydGVkXCIsXHJcbiAgICAgICAgICAgIG5lc3RlZEl0ZW1zOiBbXHJcbiAgICAgICAgICAgICAgICBcIndlbGNvbWVcIixcclxuICAgICAgICAgICAgICAgIFwiaW5zdGFsbGF0aW9uXCJcclxuICAgICAgICAgICAgXVxyXG4gICAgICAgIH0sXHJcbiAgICAgICAge1xyXG4gICAgICAgICAgICB0aXRsZTogXCJ1c2luZyBpemVcIixcclxuICAgICAgICAgICAgbmVzdGVkSXRlbXM6IFtcclxuICAgICAgICAgICAgICAgIFwiZGVwbG95aW5nIEVDUyBhcHBcIixcclxuICAgICAgICAgICAgICAgIFwiZGVwbG95aW5nIHNlcnZlcmxlc3MgYXBwXCJcclxuICAgICAgICAgICAgXVxyXG4gICAgICAgIH1cclxuICAgIF0sXHJcbiAgICBcclxufSAiXSwibmFtZXMiOlsic2lkZUJhck1lbnUiLCJtYWluTWVudSIsInRpdGxlIiwibmVzdGVkSXRlbXMiXSwic291cmNlUm9vdCI6IiJ9\n//# sourceURL=webpack-internal:///./utilities/sideBarMenu.js\n");

/***/ })

});