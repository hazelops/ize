"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
self["webpackHotUpdate_N_E"]("pages/docs",{

/***/ "./components/sideBar.js":
/*!*******************************!*\
  !*** ./components/sideBar.js ***!
  \*******************************/
/***/ (function(module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": function() { return /* binding */ SideBar; }\n/* harmony export */ });\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react/jsx-runtime */ \"./node_modules/react/jsx-runtime.js\");\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _ize__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ize */ \"./components/ize.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! react */ \"./node_modules/react/index.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(react__WEBPACK_IMPORTED_MODULE_2__);\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! next/link */ \"./node_modules/next/link.js\");\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3___default = /*#__PURE__*/__webpack_require__.n(next_link__WEBPACK_IMPORTED_MODULE_3__);\n/* harmony import */ var _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../utilities/sideBarMenu */ \"./utilities/sideBarMenu.js\");\n/* module decorator */ module = __webpack_require__.hmd(module);\n\n\n\n\n\nvar _s = $RefreshSig$();\nfunction TopElement(props) {\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        id: props.title,\n        className: \"flex items-center px-4 py-2 mt-5 text-gray-600 rounded-md hover:bg-gray-200 transition-colors duration-300 transform cursor-pointer\",\n        onClick: props.onClick,\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 8,\n            columnNumber: 9\n        },\n        __self: this,\n        children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"span\", {\n            className: \"mx-4 font-medium capitalize\",\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 9,\n                columnNumber: 13\n            },\n            __self: this,\n            children: props.title\n        })\n    }));\n}\n_c = TopElement;\nfunction NestedMenu(props) {\n    var _this = this;\n    if (props.hidden) {\n        return null;\n    }\n    var nestedList = props.nestedItems.map(function(el) {\n        var pathName = el.slice().replaceAll(\" \", \"-\");\n        var route = pathName == \"welcome\" ? \"\" : pathName;\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)((next_link__WEBPACK_IMPORTED_MODULE_3___default()), {\n            href: \"/docs/\".concat(route),\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 22,\n                columnNumber: 16\n            },\n            __self: _this,\n            children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                title: el,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 23,\n                    columnNumber: 21\n                },\n                __self: _this\n            })\n        }, el));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        className: \"flex flex-col justify-between flex-1 ml-10\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 29,\n            columnNumber: 9\n        },\n        __self: this,\n        children: nestedList\n    }));\n}\n_c1 = NestedMenu;\nfunction MenuElement(props) {\n    _s();\n    var ref = (0,react__WEBPACK_IMPORTED_MODULE_2__.useState)(false), isHidden = ref[0], setHidden = ref[1];\n    var handleClick = function handleClick() {\n        setHidden(!isHidden);\n    };\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)((react__WEBPACK_IMPORTED_MODULE_2___default().Fragment), {\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 43,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                title: props.title,\n                onClick: handleClick,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 44,\n                    columnNumber: 13\n                },\n                __self: this\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(NestedMenu, {\n                hidden: isHidden,\n                nestedItems: props.nestedItems,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 45,\n                    columnNumber: 13\n                },\n                __self: this\n            })\n        ]\n    }));\n}\n_s(MenuElement, \"Hdw5EO+DplCNBEJcNuH8tsP7WZ4=\");\n_c2 = MenuElement;\n//------------------------------------------------------------------------------------------------------------------\nfunction SideBar(props) {\n    var _this = this;\n    var mainMenu = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.mainMenu, seeAlso = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.seeAlso;\n    var menuList = mainMenu.map(function(el) {\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(MenuElement, {\n            title: el.title,\n            nestedItems: el.nestedItems,\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 57,\n                columnNumber: 13\n            },\n            __self: _this\n        }, el.title));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"div\", {\n        className: \"flex flex-col w-64 h-screen px-4 py-8 bg-white border-r\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 66,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 67,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(_ize__WEBPACK_IMPORTED_MODULE_1__[\"default\"], {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 68,\n                        columnNumber: 16\n                    },\n                    __self: this\n                })\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                className: \"flex flex-col justify-between flex-1\",\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 71,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"nav\", {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 72,\n                        columnNumber: 17\n                    },\n                    __self: this,\n                    children: [\n                        menuList,\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"hr\", {\n                            className: \"my-6 border-gray-200\",\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 74,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        }),\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(MenuElement, {\n                            title: seeAlso.title,\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 75,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        })\n                    ]\n                })\n            })\n        ]\n    }));\n};\n_c3 = SideBar;\nvar _c, _c1, _c2, _c3;\n$RefreshReg$(_c, \"TopElement\");\n$RefreshReg$(_c1, \"NestedMenu\");\n$RefreshReg$(_c2, \"MenuElement\");\n$RefreshReg$(_c3, \"SideBar\");\n\n\n;\n    var _a, _b;\n    // Legacy CSS implementations will `eval` browser code in a Node.js context\n    // to extract CSS. For backwards compatibility, we need to check we're in a\n    // browser context before continuing.\n    if (typeof self !== 'undefined' &&\n        // AMP / No-JS mode does not inject these helpers:\n        '$RefreshHelpers$' in self) {\n        var currentExports = module.__proto__.exports;\n        var prevExports = (_b = (_a = module.hot.data) === null || _a === void 0 ? void 0 : _a.prevExports) !== null && _b !== void 0 ? _b : null;\n        // This cannot happen in MainTemplate because the exports mismatch between\n        // templating and execution.\n        self.$RefreshHelpers$.registerExportsForReactRefresh(currentExports, module.id);\n        // A module can be accepted automatically based on its exports, e.g. when\n        // it is a Refresh Boundary.\n        if (self.$RefreshHelpers$.isReactRefreshBoundary(currentExports)) {\n            // Save the previous exports on update so we can compare the boundary\n            // signatures.\n            module.hot.dispose(function (data) {\n                data.prevExports = currentExports;\n            });\n            // Unconditionally accept an update to this module, we'll check if it's\n            // still a Refresh Boundary later.\n            module.hot.accept();\n            // This field is set when the previous version of this module was a\n            // Refresh Boundary, letting us know we need to check for invalidation or\n            // enqueue an update.\n            if (prevExports !== null) {\n                // A boundary can become ineligible if its exports are incompatible\n                // with the previous exports.\n                //\n                // For example, if you add/remove/change exports, we'll want to\n                // re-execute the importing modules, and force those components to\n                // re-render. Similarly, if you convert a class component to a\n                // function, we want to invalidate the boundary.\n                if (self.$RefreshHelpers$.shouldInvalidateReactRefreshBoundary(prevExports, currentExports)) {\n                    module.hot.invalidate();\n                }\n                else {\n                    self.$RefreshHelpers$.scheduleUpdate();\n                }\n            }\n        }\n        else {\n            // Since we just executed the code for the module, it's possible that the\n            // new exports made it ineligible for being a boundary.\n            // We only care about the case when we were _previously_ a boundary,\n            // because we already accepted this update (accidental side effect).\n            var isNoLongerABoundary = prevExports !== null;\n            if (isNoLongerABoundary) {\n                module.hot.invalidate();\n            }\n        }\n    }\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,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Ozs7Ozs7d0JBQ0NELFFBQVE7NkZBQ1JFLENBQUU7NEJBQUN4QixTQUFTLEVBQUMsQ0FBc0I7Ozs7Ozs7OzZGQUNuQ2MsV0FBVzs0QkFFWmYsS0FBSyxFQUFFc0IsT0FBTyxDQUFDdEIsS0FBSzs7Ozs7Ozs7Ozs7OztBQU94QyxDQUFDO01BaEN1Qm9CLE9BQU8iLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9fTl9FLy4vY29tcG9uZW50cy9zaWRlQmFyLmpzP2ZkZGEiXSwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IEl6ZSBmcm9tIFwiLi9pemVcIlxyXG5pbXBvcnQgUmVhY3QsIHsgdXNlU3RhdGUgfSBmcm9tICdyZWFjdCdcclxuaW1wb3J0IExpbmsgZnJvbSAnbmV4dC9saW5rJ1xyXG5pbXBvcnQgeyBzaWRlQmFyTWVudSB9IGZyb20gJy4uL3V0aWxpdGllcy9zaWRlQmFyTWVudSdcclxuXHJcbmZ1bmN0aW9uIFRvcEVsZW1lbnQocHJvcHMpIHtcclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPGRpdiBpZD17cHJvcHMudGl0bGV9IGNsYXNzTmFtZT1cImZsZXggaXRlbXMtY2VudGVyIHB4LTQgcHktMiBtdC01IHRleHQtZ3JheS02MDAgcm91bmRlZC1tZCBob3ZlcjpiZy1ncmF5LTIwMCB0cmFuc2l0aW9uLWNvbG9ycyBkdXJhdGlvbi0zMDAgdHJhbnNmb3JtIGN1cnNvci1wb2ludGVyXCIgb25DbGljaz17cHJvcHMub25DbGlja30+XHJcbiAgICAgICAgICAgIDxzcGFuIGNsYXNzTmFtZT1cIm14LTQgZm9udC1tZWRpdW0gY2FwaXRhbGl6ZVwiPntwcm9wcy50aXRsZX08L3NwYW4+XHJcbiAgICAgICAgPC9kaXY+XHJcbiAgICApXHJcbn1cclxuXHJcbmZ1bmN0aW9uIE5lc3RlZE1lbnUocHJvcHMpIHtcclxuICAgIGlmIChwcm9wcy5oaWRkZW4pIHtcclxuICAgICAgICByZXR1cm4gbnVsbFxyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IG5lc3RlZExpc3QgPSBwcm9wcy5uZXN0ZWRJdGVtcy5tYXAoZWwgPT4ge1xyXG4gICAgICAgIGNvbnN0IHBhdGhOYW1lID0gZWwuc2xpY2UoKS5yZXBsYWNlQWxsKFwiIFwiLCBcIi1cIilcclxuICAgICAgICBsZXQgcm91dGUgPSBwYXRoTmFtZSA9PSBcIndlbGNvbWVcIj8gXCJcIiA6IHBhdGhOYW1lXHJcbiAgICAgICAgcmV0dXJuIDxMaW5rIGtleT17ZWx9IGhyZWY9e2AvZG9jcy8ke3JvdXRlfWB9PlxyXG4gICAgICAgICAgICAgICAgICAgIDxUb3BFbGVtZW50IHRpdGxlPXtlbH0vPlxyXG4gICAgICAgICAgICAgICAgPC9MaW5rPlxyXG4gICAgICAgIFxyXG4gICAgfSlcclxuXHJcbiAgICByZXR1cm4gKFxyXG4gICAgICAgIDxkaXYgY2xhc3NOYW1lPVwiZmxleCBmbGV4LWNvbCBqdXN0aWZ5LWJldHdlZW4gZmxleC0xIG1sLTEwXCI+XHJcbiAgICAgICAgICAgIHtuZXN0ZWRMaXN0fVxyXG4gICAgICAgIDwvZGl2PlxyXG4gICAgKVxyXG59XHJcblxyXG5mdW5jdGlvbiBNZW51RWxlbWVudChwcm9wcykge1xyXG4gICAgY29uc3QgW2lzSGlkZGVuLCBzZXRIaWRkZW5dID0gdXNlU3RhdGUoZmFsc2UpXHJcblxyXG4gICAgY29uc3QgaGFuZGxlQ2xpY2sgPSBmdW5jdGlvbigpIHtcclxuICAgICAgICBzZXRIaWRkZW4oIWlzSGlkZGVuKVxyXG4gICAgfSBcclxuXHJcbiAgICByZXR1cm4gKFxyXG4gICAgICAgIDxSZWFjdC5GcmFnbWVudD5cclxuICAgICAgICAgICAgPFRvcEVsZW1lbnQgdGl0bGU9e3Byb3BzLnRpdGxlfSBvbkNsaWNrPXtoYW5kbGVDbGlja30gLz5cclxuICAgICAgICAgICAgPE5lc3RlZE1lbnUgaGlkZGVuPXtpc0hpZGRlbn0gbmVzdGVkSXRlbXM9e3Byb3BzLm5lc3RlZEl0ZW1zfSAvPlxyXG4gICAgICAgIDwvUmVhY3QuRnJhZ21lbnQ+XHJcbiAgICAgICAgXHJcbiAgICApXHJcbn1cclxuLy8tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS1cclxuXHJcbmV4cG9ydCBkZWZhdWx0IGZ1bmN0aW9uIFNpZGVCYXIocHJvcHMpIHtcclxuICAgIGNvbnN0IHsgbWFpbk1lbnUsIHNlZUFsc28gfSA9IHNpZGVCYXJNZW51XHJcblxyXG4gICAgY29uc3QgbWVudUxpc3QgPSBtYWluTWVudS5tYXAoZWwgPT4ge1xyXG4gICAgICAgIHJldHVybiAoXHJcbiAgICAgICAgICAgIDxNZW51RWxlbWVudFxyXG4gICAgICAgICAgICAgICAga2V5PXtlbC50aXRsZX1cclxuICAgICAgICAgICAgICAgIHRpdGxlPXtlbC50aXRsZX1cclxuICAgICAgICAgICAgICAgIG5lc3RlZEl0ZW1zPXtlbC5uZXN0ZWRJdGVtc31cclxuICAgICAgICAgICAgIC8+XHJcbiAgICAgICAgKVxyXG4gICAgfSlcclxuXHJcbiAgICByZXR1cm4gKFxyXG4gICAgICAgIDxkaXYgY2xhc3NOYW1lPVwiZmxleCBmbGV4LWNvbCB3LTY0IGgtc2NyZWVuIHB4LTQgcHktOCBiZy13aGl0ZSBib3JkZXItclwiPlxyXG4gICAgICAgICAgICA8ZGl2PlxyXG4gICAgICAgICAgICAgICA8SXplIC8+IFxyXG4gICAgICAgICAgICA8L2Rpdj5cclxuXHJcbiAgICAgICAgICAgIDxkaXYgY2xhc3NOYW1lPVwiZmxleCBmbGV4LWNvbCBqdXN0aWZ5LWJldHdlZW4gZmxleC0xXCI+XHJcbiAgICAgICAgICAgICAgICA8bmF2PlxyXG4gICAgICAgICAgICAgICAgICAgIHttZW51TGlzdH1cclxuICAgICAgICAgICAgICAgICAgICA8aHIgY2xhc3NOYW1lPVwibXktNiBib3JkZXItZ3JheS0yMDBcIiAvPlxyXG4gICAgICAgICAgICAgICAgICAgIDxNZW51RWxlbWVudFxyXG5cclxuICAgICAgICAgICAgICAgICAgICB0aXRsZT17c2VlQWxzby50aXRsZX1cclxuICAgICAgICAgICAgICAgICAgICAgICAgXHJcbiAgICAgICAgICAgICAgICAgICAgLz5cclxuICAgICAgICAgICAgICAgIDwvbmF2PlxyXG4gICAgICAgICAgICA8L2Rpdj5cclxuICAgICAgICA8L2Rpdj5cclxuICAgIClcclxufSJdLCJuYW1lcyI6WyJJemUiLCJSZWFjdCIsInVzZVN0YXRlIiwiTGluayIsInNpZGVCYXJNZW51IiwiVG9wRWxlbWVudCIsInByb3BzIiwiZGl2IiwiaWQiLCJ0aXRsZSIsImNsYXNzTmFtZSIsIm9uQ2xpY2siLCJzcGFuIiwiTmVzdGVkTWVudSIsImhpZGRlbiIsIm5lc3RlZExpc3QiLCJuZXN0ZWRJdGVtcyIsIm1hcCIsImVsIiwicGF0aE5hbWUiLCJzbGljZSIsInJlcGxhY2VBbGwiLCJyb3V0ZSIsImhyZWYiLCJNZW51RWxlbWVudCIsImlzSGlkZGVuIiwic2V0SGlkZGVuIiwiaGFuZGxlQ2xpY2siLCJGcmFnbWVudCIsIlNpZGVCYXIiLCJtYWluTWVudSIsInNlZUFsc28iLCJtZW51TGlzdCIsIm5hdiIsImhyIl0sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./components/sideBar.js\n");

/***/ })

});