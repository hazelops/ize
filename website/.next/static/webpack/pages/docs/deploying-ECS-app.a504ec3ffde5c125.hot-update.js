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

/***/ "./components/sideBar.js":
/*!*******************************!*\
  !*** ./components/sideBar.js ***!
  \*******************************/
/***/ (function(module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": function() { return /* binding */ SideBar; }\n/* harmony export */ });\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react/jsx-runtime */ \"./node_modules/react/jsx-runtime.js\");\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _ize__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ize */ \"./components/ize.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! react */ \"./node_modules/react/index.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(react__WEBPACK_IMPORTED_MODULE_2__);\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! next/link */ \"./node_modules/next/link.js\");\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3___default = /*#__PURE__*/__webpack_require__.n(next_link__WEBPACK_IMPORTED_MODULE_3__);\n/* harmony import */ var _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../utilities/sideBarMenu */ \"./utilities/sideBarMenu.js\");\n/* module decorator */ module = __webpack_require__.hmd(module);\n\n\n\n\n\nvar _s = $RefreshSig$();\nfunction TopElement(props) {\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        id: props.title,\n        className: \"flex items-center px-4 py-2 mt-5 text-gray-600 rounded-md hover:bg-gray-200 transition-colors duration-300 transform cursor-pointer\",\n        onClick: props.onClick,\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 8,\n            columnNumber: 9\n        },\n        __self: this,\n        children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"span\", {\n            className: \"mx-4 font-medium capitalize\",\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 9,\n                columnNumber: 13\n            },\n            __self: this,\n            children: props.title\n        })\n    }));\n}\n_c = TopElement;\nfunction NestedMenu(props) {\n    var _this = this;\n    if (props.hidden) {\n        return null;\n    }\n    var nestedList = props.nestedItems.map(function(el) {\n        var pathName = el.slice().replaceAll(\" \", \"-\");\n        var route = pathName == \"welcome\" ? \"\" : pathName;\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)((next_link__WEBPACK_IMPORTED_MODULE_3___default()), {\n            href: \"/docs/\".concat(route),\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 22,\n                columnNumber: 16\n            },\n            __self: _this,\n            children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                title: el,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 23,\n                    columnNumber: 21\n                },\n                __self: _this\n            })\n        }, el));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        className: \"flex flex-col justify-between flex-1 ml-10\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 29,\n            columnNumber: 9\n        },\n        __self: this,\n        children: nestedList\n    }));\n}\n_c1 = NestedMenu;\nfunction MenuElement(props) {\n    _s();\n    var ref = (0,react__WEBPACK_IMPORTED_MODULE_2__.useState)(false), isHidden = ref[0], setHidden = ref[1];\n    var handleClick = function handleClick() {\n        setHidden(!isHidden);\n    };\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)((react__WEBPACK_IMPORTED_MODULE_2___default().Fragment), {\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 43,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                title: props.title,\n                onClick: handleClick,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 44,\n                    columnNumber: 13\n                },\n                __self: this\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(NestedMenu, {\n                hidden: isHidden,\n                nestedItems: props.nestedItems,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 45,\n                    columnNumber: 13\n                },\n                __self: this\n            })\n        ]\n    }));\n}\n_s(MenuElement, \"Hdw5EO+DplCNBEJcNuH8tsP7WZ4=\");\n_c2 = MenuElement;\n//------------------------------------------------------------------------------------------------------------------\nfunction SideBar(props) {\n    var _this = this;\n    var mainMenu = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.mainMenu, seeAlso = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.seeAlso;\n    var menuList = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.mainMenu.map(function(el) {\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(MenuElement, {\n            title: el.title,\n            nestedItems: el.nestedItems,\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 57,\n                columnNumber: 13\n            },\n            __self: _this\n        }, el.title));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"div\", {\n        className: \"flex flex-col w-64 h-screen px-4 py-8 bg-white border-r\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 66,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 67,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(_ize__WEBPACK_IMPORTED_MODULE_1__[\"default\"], {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 68,\n                        columnNumber: 16\n                    },\n                    __self: this\n                })\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                className: \"flex flex-col justify-between flex-1\",\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 71,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"nav\", {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 72,\n                        columnNumber: 17\n                    },\n                    __self: this,\n                    children: [\n                        menuList,\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"hr\", {\n                            className: \"my-6 border-gray-200\",\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 74,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        }),\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                            title: \"See Also\",\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 75,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        })\n                    ]\n                })\n            })\n        ]\n    }));\n};\n_c3 = SideBar;\nvar _c, _c1, _c2, _c3;\n$RefreshReg$(_c, \"TopElement\");\n$RefreshReg$(_c1, \"NestedMenu\");\n$RefreshReg$(_c2, \"MenuElement\");\n$RefreshReg$(_c3, \"SideBar\");\n\n\n;\n    var _a, _b;\n    // Legacy CSS implementations will `eval` browser code in a Node.js context\n    // to extract CSS. For backwards compatibility, we need to check we're in a\n    // browser context before continuing.\n    if (typeof self !== 'undefined' &&\n        // AMP / No-JS mode does not inject these helpers:\n        '$RefreshHelpers$' in self) {\n        var currentExports = module.__proto__.exports;\n        var prevExports = (_b = (_a = module.hot.data) === null || _a === void 0 ? void 0 : _a.prevExports) !== null && _b !== void 0 ? _b : null;\n        // This cannot happen in MainTemplate because the exports mismatch between\n        // templating and execution.\n        self.$RefreshHelpers$.registerExportsForReactRefresh(currentExports, module.id);\n        // A module can be accepted automatically based on its exports, e.g. when\n        // it is a Refresh Boundary.\n        if (self.$RefreshHelpers$.isReactRefreshBoundary(currentExports)) {\n            // Save the previous exports on update so we can compare the boundary\n            // signatures.\n            module.hot.dispose(function (data) {\n                data.prevExports = currentExports;\n            });\n            // Unconditionally accept an update to this module, we'll check if it's\n            // still a Refresh Boundary later.\n            module.hot.accept();\n            // This field is set when the previous version of this module was a\n            // Refresh Boundary, letting us know we need to check for invalidation or\n            // enqueue an update.\n            if (prevExports !== null) {\n                // A boundary can become ineligible if its exports are incompatible\n                // with the previous exports.\n                //\n                // For example, if you add/remove/change exports, we'll want to\n                // re-execute the importing modules, and force those components to\n                // re-render. Similarly, if you convert a class component to a\n                // function, we want to invalidate the boundary.\n                if (self.$RefreshHelpers$.shouldInvalidateReactRefreshBoundary(prevExports, currentExports)) {\n                    module.hot.invalidate();\n                }\n                else {\n                    self.$RefreshHelpers$.scheduleUpdate();\n                }\n            }\n        }\n        else {\n            // Since we just executed the code for the module, it's possible that the\n            // new exports made it ineligible for being a boundary.\n            // We only care about the case when we were _previously_ a boundary,\n            // because we already accepted this update (accidental side effect).\n            var isNoLongerABoundary = prevExports !== null;\n            if (isNoLongerABoundary) {\n                module.hot.invalidate();\n            }\n        }\n    }\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,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QUNDRCxRQUFROzZGQUNSRSxDQUFFOzRCQUFDeEIsU0FBUyxFQUFDLENBQXNCOzs7Ozs7Ozs2RkFDbkNMLFVBQVU7NEJBQUNJLEtBQUssRUFBQyxDQUFVOzs7Ozs7Ozs7Ozs7O0FBS2hELENBQUM7TUE1QnVCb0IsT0FBTyIsInNvdXJjZXMiOlsid2VicGFjazovL19OX0UvLi9jb21wb25lbnRzL3NpZGVCYXIuanM/ZmRkYSJdLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgSXplIGZyb20gXCIuL2l6ZVwiXHJcbmltcG9ydCBSZWFjdCwgeyB1c2VTdGF0ZSB9IGZyb20gJ3JlYWN0J1xyXG5pbXBvcnQgTGluayBmcm9tICduZXh0L2xpbmsnXHJcbmltcG9ydCB7IHNpZGVCYXJNZW51IH0gZnJvbSAnLi4vdXRpbGl0aWVzL3NpZGVCYXJNZW51J1xyXG5cclxuZnVuY3Rpb24gVG9wRWxlbWVudChwcm9wcykge1xyXG4gICAgcmV0dXJuIChcclxuICAgICAgICA8ZGl2IGlkPXtwcm9wcy50aXRsZX0gY2xhc3NOYW1lPVwiZmxleCBpdGVtcy1jZW50ZXIgcHgtNCBweS0yIG10LTUgdGV4dC1ncmF5LTYwMCByb3VuZGVkLW1kIGhvdmVyOmJnLWdyYXktMjAwIHRyYW5zaXRpb24tY29sb3JzIGR1cmF0aW9uLTMwMCB0cmFuc2Zvcm0gY3Vyc29yLXBvaW50ZXJcIiBvbkNsaWNrPXtwcm9wcy5vbkNsaWNrfT5cclxuICAgICAgICAgICAgPHNwYW4gY2xhc3NOYW1lPVwibXgtNCBmb250LW1lZGl1bSBjYXBpdGFsaXplXCI+e3Byb3BzLnRpdGxlfTwvc3Bhbj5cclxuICAgICAgICA8L2Rpdj5cclxuICAgIClcclxufVxyXG5cclxuZnVuY3Rpb24gTmVzdGVkTWVudShwcm9wcykge1xyXG4gICAgaWYgKHByb3BzLmhpZGRlbikge1xyXG4gICAgICAgIHJldHVybiBudWxsXHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgbmVzdGVkTGlzdCA9IHByb3BzLm5lc3RlZEl0ZW1zLm1hcChlbCA9PiB7XHJcbiAgICAgICAgY29uc3QgcGF0aE5hbWUgPSBlbC5zbGljZSgpLnJlcGxhY2VBbGwoXCIgXCIsIFwiLVwiKVxyXG4gICAgICAgIGxldCByb3V0ZSA9IHBhdGhOYW1lID09IFwid2VsY29tZVwiPyBcIlwiIDogcGF0aE5hbWVcclxuICAgICAgICByZXR1cm4gPExpbmsga2V5PXtlbH0gaHJlZj17YC9kb2NzLyR7cm91dGV9YH0+XHJcbiAgICAgICAgICAgICAgICAgICAgPFRvcEVsZW1lbnQgdGl0bGU9e2VsfS8+XHJcbiAgICAgICAgICAgICAgICA8L0xpbms+XHJcbiAgICAgICAgXHJcbiAgICB9KVxyXG5cclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPGRpdiBjbGFzc05hbWU9XCJmbGV4IGZsZXgtY29sIGp1c3RpZnktYmV0d2VlbiBmbGV4LTEgbWwtMTBcIj5cclxuICAgICAgICAgICAge25lc3RlZExpc3R9XHJcbiAgICAgICAgPC9kaXY+XHJcbiAgICApXHJcbn1cclxuXHJcbmZ1bmN0aW9uIE1lbnVFbGVtZW50KHByb3BzKSB7XHJcbiAgICBjb25zdCBbaXNIaWRkZW4sIHNldEhpZGRlbl0gPSB1c2VTdGF0ZShmYWxzZSlcclxuXHJcbiAgICBjb25zdCBoYW5kbGVDbGljayA9IGZ1bmN0aW9uKCkge1xyXG4gICAgICAgIHNldEhpZGRlbighaXNIaWRkZW4pXHJcbiAgICB9IFxyXG5cclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPFJlYWN0LkZyYWdtZW50PlxyXG4gICAgICAgICAgICA8VG9wRWxlbWVudCB0aXRsZT17cHJvcHMudGl0bGV9IG9uQ2xpY2s9e2hhbmRsZUNsaWNrfSAvPlxyXG4gICAgICAgICAgICA8TmVzdGVkTWVudSBoaWRkZW49e2lzSGlkZGVufSBuZXN0ZWRJdGVtcz17cHJvcHMubmVzdGVkSXRlbXN9IC8+XHJcbiAgICAgICAgPC9SZWFjdC5GcmFnbWVudD5cclxuICAgICAgICBcclxuICAgIClcclxufVxyXG4vLy0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLVxyXG5cclxuZXhwb3J0IGRlZmF1bHQgZnVuY3Rpb24gU2lkZUJhcihwcm9wcykge1xyXG4gICAgY29uc3QgeyBtYWluTWVudSwgc2VlQWxzbyB9ID0gc2lkZUJhck1lbnVcclxuXHJcbiAgICBjb25zdCBtZW51TGlzdCA9IHNpZGVCYXJNZW51Lm1haW5NZW51Lm1hcChlbCA9PiB7XHJcbiAgICAgICAgcmV0dXJuIChcclxuICAgICAgICAgICAgPE1lbnVFbGVtZW50XHJcbiAgICAgICAgICAgICAgICBrZXk9e2VsLnRpdGxlfVxyXG4gICAgICAgICAgICAgICAgdGl0bGU9e2VsLnRpdGxlfVxyXG4gICAgICAgICAgICAgICAgbmVzdGVkSXRlbXM9e2VsLm5lc3RlZEl0ZW1zfVxyXG4gICAgICAgICAgICAgLz5cclxuICAgICAgICApXHJcbiAgICB9KVxyXG5cclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPGRpdiBjbGFzc05hbWU9XCJmbGV4IGZsZXgtY29sIHctNjQgaC1zY3JlZW4gcHgtNCBweS04IGJnLXdoaXRlIGJvcmRlci1yXCI+XHJcbiAgICAgICAgICAgIDxkaXY+XHJcbiAgICAgICAgICAgICAgIDxJemUgLz4gXHJcbiAgICAgICAgICAgIDwvZGl2PlxyXG5cclxuICAgICAgICAgICAgPGRpdiBjbGFzc05hbWU9XCJmbGV4IGZsZXgtY29sIGp1c3RpZnktYmV0d2VlbiBmbGV4LTFcIj5cclxuICAgICAgICAgICAgICAgIDxuYXY+XHJcbiAgICAgICAgICAgICAgICAgICAge21lbnVMaXN0fVxyXG4gICAgICAgICAgICAgICAgICAgIDxociBjbGFzc05hbWU9XCJteS02IGJvcmRlci1ncmF5LTIwMFwiIC8+XHJcbiAgICAgICAgICAgICAgICAgICAgPFRvcEVsZW1lbnQgdGl0bGU9XCJTZWUgQWxzb1wiIC8+XHJcbiAgICAgICAgICAgICAgICA8L25hdj5cclxuICAgICAgICAgICAgPC9kaXY+XHJcbiAgICAgICAgPC9kaXY+XHJcbiAgICApXHJcbn0iXSwibmFtZXMiOlsiSXplIiwiUmVhY3QiLCJ1c2VTdGF0ZSIsIkxpbmsiLCJzaWRlQmFyTWVudSIsIlRvcEVsZW1lbnQiLCJwcm9wcyIsImRpdiIsImlkIiwidGl0bGUiLCJjbGFzc05hbWUiLCJvbkNsaWNrIiwic3BhbiIsIk5lc3RlZE1lbnUiLCJoaWRkZW4iLCJuZXN0ZWRMaXN0IiwibmVzdGVkSXRlbXMiLCJtYXAiLCJlbCIsInBhdGhOYW1lIiwic2xpY2UiLCJyZXBsYWNlQWxsIiwicm91dGUiLCJocmVmIiwiTWVudUVsZW1lbnQiLCJpc0hpZGRlbiIsInNldEhpZGRlbiIsImhhbmRsZUNsaWNrIiwiRnJhZ21lbnQiLCJTaWRlQmFyIiwibWFpbk1lbnUiLCJzZWVBbHNvIiwibWVudUxpc3QiLCJuYXYiLCJociJdLCJzb3VyY2VSb290IjoiIn0=\n//# sourceURL=webpack-internal:///./components/sideBar.js\n");

/***/ })

});