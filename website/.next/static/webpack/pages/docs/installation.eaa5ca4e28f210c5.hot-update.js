"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
self["webpackHotUpdate_N_E"]("pages/docs/installation",{

/***/ "./components/sideBar.js":
/*!*******************************!*\
  !*** ./components/sideBar.js ***!
  \*******************************/
/***/ (function(module, __webpack_exports__, __webpack_require__) {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": function() { return /* binding */ SideBar; }\n/* harmony export */ });\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react/jsx-runtime */ \"./node_modules/react/jsx-runtime.js\");\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _ize__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ize */ \"./components/ize.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! react */ \"./node_modules/react/index.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(react__WEBPACK_IMPORTED_MODULE_2__);\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! next/link */ \"./node_modules/next/link.js\");\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3___default = /*#__PURE__*/__webpack_require__.n(next_link__WEBPACK_IMPORTED_MODULE_3__);\n/* harmony import */ var _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../utilities/sideBarMenu */ \"./utilities/sideBarMenu.js\");\n/* module decorator */ module = __webpack_require__.hmd(module);\n\n\n\n\n\nvar _s = $RefreshSig$();\nfunction TopElement(props) {\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        id: props.title,\n        className: \"flex items-center px-4 py-2 mt-5 text-gray-600 rounded-md hover:bg-gray-200 transition-colors duration-300 transform cursor-pointer\",\n        onClick: props.onClick,\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 8,\n            columnNumber: 13\n        },\n        __self: this,\n        children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"span\", {\n            className: \"mx-4 font-medium capitalize\",\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 9,\n                columnNumber: 17\n            },\n            __self: this,\n            children: props.title\n        })\n    }));\n}\n_c = TopElement;\nfunction NestedMenu(props) {\n    var _this = this;\n    if (props.hidden) {\n        return null;\n    }\n    var nestedList = props.nestedItems.map(function(el) {\n        var pathName = el.slice().replaceAll(\" \", \"-\");\n        var route = pathName == \"welcome\" ? \"\" : pathName;\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)((next_link__WEBPACK_IMPORTED_MODULE_3___default()), {\n            href: \"/docs/\".concat(route),\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 22,\n                columnNumber: 16\n            },\n            __self: _this,\n            children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"a\", {\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 25,\n                    columnNumber: 21\n                },\n                __self: _this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                    title: el,\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 25,\n                        columnNumber: 24\n                    },\n                    __self: _this\n                })\n            })\n        }, props.nestedItems.indexOf(el)));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        className: \"flex flex-col justify-between flex-1 ml-10\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 30,\n            columnNumber: 9\n        },\n        __self: this,\n        children: nestedList\n    }));\n}\n_c1 = NestedMenu;\nfunction MenuElement(props) {\n    _s();\n    var ref = (0,react__WEBPACK_IMPORTED_MODULE_2__.useState)(false), isHidden = ref[0], setHidden = ref[1];\n    var handleClick = function handleClick() {\n        setHidden(!isHidden);\n    };\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)((react__WEBPACK_IMPORTED_MODULE_2___default().Fragment), {\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 44,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                title: props.title,\n                onClick: handleClick,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 45,\n                    columnNumber: 13\n                },\n                __self: this\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(NestedMenu, {\n                hidden: isHidden,\n                nestedItems: props.nestedItems,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 46,\n                    columnNumber: 13\n                },\n                __self: this\n            })\n        ]\n    }));\n}\n_s(MenuElement, \"Hdw5EO+DplCNBEJcNuH8tsP7WZ4=\");\n_c2 = MenuElement;\n//------------------------------------------------------------------------------------------------------------------\nfunction SideBar(props) {\n    var _this = this;\n    var mainMenu = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.mainMenu, seeAlso = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.seeAlso;\n    var menuList = mainMenu.map(function(el) {\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(MenuElement, {\n            title: el.title,\n            nestedItems: el.nestedItems,\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 58,\n                columnNumber: 13\n            },\n            __self: _this\n        }, mainMenu.indexOf(el)));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"div\", {\n        className: \"flex flex-col w-64 h-screen px-4 py-8 bg-white border-r\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 67,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 68,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(_ize__WEBPACK_IMPORTED_MODULE_1__[\"default\"], {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 69,\n                        columnNumber: 16\n                    },\n                    __self: this\n                })\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                className: \"flex flex-col justify-between flex-1\",\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 72,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"nav\", {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 73,\n                        columnNumber: 17\n                    },\n                    __self: this,\n                    children: [\n                        menuList,\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"hr\", {\n                            className: \"my-6 border-gray-200\",\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 75,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        }),\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                            title: seeAlso.title,\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 76,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        })\n                    ]\n                })\n            })\n        ]\n    }));\n};\n_c3 = SideBar;\nvar _c, _c1, _c2, _c3;\n$RefreshReg$(_c, \"TopElement\");\n$RefreshReg$(_c1, \"NestedMenu\");\n$RefreshReg$(_c2, \"MenuElement\");\n$RefreshReg$(_c3, \"SideBar\");\n\n\n;\n    var _a, _b;\n    // Legacy CSS implementations will `eval` browser code in a Node.js context\n    // to extract CSS. For backwards compatibility, we need to check we're in a\n    // browser context before continuing.\n    if (typeof self !== 'undefined' &&\n        // AMP / No-JS mode does not inject these helpers:\n        '$RefreshHelpers$' in self) {\n        var currentExports = module.__proto__.exports;\n        var prevExports = (_b = (_a = module.hot.data) === null || _a === void 0 ? void 0 : _a.prevExports) !== null && _b !== void 0 ? _b : null;\n        // This cannot happen in MainTemplate because the exports mismatch between\n        // templating and execution.\n        self.$RefreshHelpers$.registerExportsForReactRefresh(currentExports, module.id);\n        // A module can be accepted automatically based on its exports, e.g. when\n        // it is a Refresh Boundary.\n        if (self.$RefreshHelpers$.isReactRefreshBoundary(currentExports)) {\n            // Save the previous exports on update so we can compare the boundary\n            // signatures.\n            module.hot.dispose(function (data) {\n                data.prevExports = currentExports;\n            });\n            // Unconditionally accept an update to this module, we'll check if it's\n            // still a Refresh Boundary later.\n            module.hot.accept();\n            // This field is set when the previous version of this module was a\n            // Refresh Boundary, letting us know we need to check for invalidation or\n            // enqueue an update.\n            if (prevExports !== null) {\n                // A boundary can become ineligible if its exports are incompatible\n                // with the previous exports.\n                //\n                // For example, if you add/remove/change exports, we'll want to\n                // re-execute the importing modules, and force those components to\n                // re-render. Similarly, if you convert a class component to a\n                // function, we want to invalidate the boundary.\n                if (self.$RefreshHelpers$.shouldInvalidateReactRefreshBoundary(prevExports, currentExports)) {\n                    module.hot.invalidate();\n                }\n                else {\n                    self.$RefreshHelpers$.scheduleUpdate();\n                }\n            }\n        }\n        else {\n            // Since we just executed the code for the module, it's possible that the\n            // new exports made it ineligible for being a boundary.\n            // We only care about the case when we were _previously_ a boundary,\n            // because we already accepted this update (accidental side effect).\n            var isNoLongerABoundary = prevExports !== null;\n            if (isNoLongerABoundary) {\n                module.hot.invalidate();\n            }\n        }\n    }\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiLi9jb21wb25lbnRzL3NpZGVCYXIuanMuanMiLCJtYXBwaW5ncyI6Ijs7Ozs7Ozs7Ozs7Ozs7QUFBdUI7QUFDZ0I7QUFDWDtBQUMwQjs7U0FFN0NLLFVBQVUsQ0FBQ0MsS0FBSyxFQUFFLENBQUM7SUFDeEIsTUFBTSxzRUFDR0MsQ0FBRztRQUFDQyxFQUFFLEVBQUVGLEtBQUssQ0FBQ0csS0FBSztRQUFFQyxTQUFTLEVBQUMsQ0FBcUk7UUFBQ0MsT0FBTyxFQUFFTCxLQUFLLENBQUNLLE9BQU87Ozs7Ozs7dUZBQ3ZMQyxDQUFJO1lBQUNGLFNBQVMsRUFBQyxDQUE2Qjs7Ozs7OztzQkFBRUosS0FBSyxDQUFDRyxLQUFLOzs7QUFHMUUsQ0FBQztLQU5RSixVQUFVO1NBUVZRLFVBQVUsQ0FBQ1AsS0FBSyxFQUFFLENBQUM7O0lBQ3hCLEVBQUUsRUFBRUEsS0FBSyxDQUFDUSxNQUFNLEVBQUUsQ0FBQztRQUNmLE1BQU0sQ0FBQyxJQUFJO0lBQ2YsQ0FBQztJQUVELEdBQUssQ0FBQ0MsVUFBVSxHQUFHVCxLQUFLLENBQUNVLFdBQVcsQ0FBQ0MsR0FBRyxDQUFDQyxRQUFRLENBQVJBLEVBQUUsRUFBSSxDQUFDO1FBQzVDLEdBQUssQ0FBQ0MsUUFBUSxHQUFHRCxFQUFFLENBQUNFLEtBQUssR0FBR0MsVUFBVSxDQUFDLENBQUcsSUFBRSxDQUFHO1FBQy9DLEdBQUcsQ0FBQ0MsS0FBSyxHQUFHSCxRQUFRLElBQUksQ0FBUyxXQUFFLENBQUUsSUFBR0EsUUFBUTtRQUNoRCxNQUFNLHNFQUFFaEIsa0RBQUk7WUFFSm9CLElBQUksRUFBRyxDQUFNLFFBQVEsT0FBTkQsS0FBSzs7Ozs7OzsyRkFDZkUsQ0FBQzs7Ozs7OzsrRkFBRW5CLFVBQVU7b0JBQUNJLEtBQUssRUFBRVMsRUFBRTs7Ozs7Ozs7O1dBRnZCWixLQUFLLENBQUNVLFdBQVcsQ0FBQ1MsT0FBTyxDQUFDUCxFQUFFO0lBSTdDLENBQUM7SUFFRCxNQUFNLHNFQUNEWCxDQUFHO1FBQUNHLFNBQVMsRUFBQyxDQUE0Qzs7Ozs7OztrQkFDdERLLFVBQVU7O0FBR3ZCLENBQUM7TUFwQlFGLFVBQVU7U0FzQlZhLFdBQVcsQ0FBQ3BCLEtBQUssRUFBRSxDQUFDOztJQUN6QixHQUFLLENBQXlCSixHQUFlLEdBQWZBLCtDQUFRLENBQUMsS0FBSyxHQUFyQ3lCLFFBQVEsR0FBZXpCLEdBQWUsS0FBNUIwQixTQUFTLEdBQUkxQixHQUFlO0lBRTdDLEdBQUssQ0FBQzJCLFdBQVcsR0FBRyxRQUFRLENBQXRCQSxXQUFXLEdBQWMsQ0FBQztRQUM1QkQsU0FBUyxFQUFFRCxRQUFRO0lBQ3ZCLENBQUM7SUFFRCxNQUFNLHVFQUNEMUIsdURBQWM7Ozs7Ozs7O2lGQUNWSSxVQUFVO2dCQUFDSSxLQUFLLEVBQUVILEtBQUssQ0FBQ0csS0FBSztnQkFBRUUsT0FBTyxFQUFFa0IsV0FBVzs7Ozs7Ozs7aUZBQ25EaEIsVUFBVTtnQkFBQ0MsTUFBTSxFQUFFYSxRQUFRO2dCQUFFWCxXQUFXLEVBQUVWLEtBQUssQ0FBQ1UsV0FBVzs7Ozs7Ozs7OztBQUl4RSxDQUFDO0dBZFFVLFdBQVc7TUFBWEEsV0FBVztBQWVwQixFQUFvSDtBQUVyRyxRQUFRLENBQUNLLE9BQU8sQ0FBQ3pCLEtBQUssRUFBRSxDQUFDOztJQUNwQyxHQUFLLENBQUcwQixRQUFRLEdBQWM1Qix3RUFBZCxFQUFFNkIsT0FBTyxHQUFLN0IsdUVBQUw7SUFFekIsR0FBSyxDQUFDOEIsUUFBUSxHQUFHRixRQUFRLENBQUNmLEdBQUcsQ0FBQ0MsUUFBUSxDQUFSQSxFQUFFLEVBQUksQ0FBQztRQUNqQyxNQUFNLHNFQUNEUSxXQUFXO1lBRVJqQixLQUFLLEVBQUVTLEVBQUUsQ0FBQ1QsS0FBSztZQUNmTyxXQUFXLEVBQUVFLEVBQUUsQ0FBQ0YsV0FBVzs7Ozs7OztXQUZ0QmdCLFFBQVEsQ0FBQ1AsT0FBTyxDQUFDUCxFQUFFO0lBS3BDLENBQUM7SUFFRCxNQUFNLHVFQUNEWCxDQUFHO1FBQUNHLFNBQVMsRUFBQyxDQUF5RDs7Ozs7Ozs7aUZBQ25FSCxDQUFHOzs7Ozs7OytGQUNBUCw0Q0FBRzs7Ozs7Ozs7O2lGQUdOTyxDQUFHO2dCQUFDRyxTQUFTLEVBQUMsQ0FBc0M7Ozs7Ozs7Z0dBQ2hEeUIsQ0FBRzs7Ozs7Ozs7d0JBQ0NELFFBQVE7NkZBQ1JFLENBQUU7NEJBQUMxQixTQUFTLEVBQUMsQ0FBc0I7Ozs7Ozs7OzZGQUNuQ0wsVUFBVTs0QkFDUEksS0FBSyxFQUFFd0IsT0FBTyxDQUFDeEIsS0FBSzs7Ozs7Ozs7Ozs7OztBQU01QyxDQUFDO01BOUJ1QnNCLE9BQU8iLCJzb3VyY2VzIjpbIndlYnBhY2s6Ly9fTl9FLy4vY29tcG9uZW50cy9zaWRlQmFyLmpzP2ZkZGEiXSwic291cmNlc0NvbnRlbnQiOlsiaW1wb3J0IEl6ZSBmcm9tIFwiLi9pemVcIlxyXG5pbXBvcnQgUmVhY3QsIHsgdXNlU3RhdGUgfSBmcm9tICdyZWFjdCdcclxuaW1wb3J0IExpbmsgZnJvbSAnbmV4dC9saW5rJ1xyXG5pbXBvcnQgeyBzaWRlQmFyTWVudSB9IGZyb20gJy4uL3V0aWxpdGllcy9zaWRlQmFyTWVudSdcclxuXHJcbmZ1bmN0aW9uIFRvcEVsZW1lbnQocHJvcHMpIHtcclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgICAgIDxkaXYgaWQ9e3Byb3BzLnRpdGxlfSBjbGFzc05hbWU9XCJmbGV4IGl0ZW1zLWNlbnRlciBweC00IHB5LTIgbXQtNSB0ZXh0LWdyYXktNjAwIHJvdW5kZWQtbWQgaG92ZXI6YmctZ3JheS0yMDAgdHJhbnNpdGlvbi1jb2xvcnMgZHVyYXRpb24tMzAwIHRyYW5zZm9ybSBjdXJzb3ItcG9pbnRlclwiIG9uQ2xpY2s9e3Byb3BzLm9uQ2xpY2t9PlxyXG4gICAgICAgICAgICAgICAgPHNwYW4gY2xhc3NOYW1lPVwibXgtNCBmb250LW1lZGl1bSBjYXBpdGFsaXplXCI+e3Byb3BzLnRpdGxlfTwvc3Bhbj5cclxuICAgICAgICAgICAgPC9kaXY+XHJcbiAgICApXHJcbn1cclxuXHJcbmZ1bmN0aW9uIE5lc3RlZE1lbnUocHJvcHMpIHtcclxuICAgIGlmIChwcm9wcy5oaWRkZW4pIHtcclxuICAgICAgICByZXR1cm4gbnVsbFxyXG4gICAgfVxyXG5cclxuICAgIGNvbnN0IG5lc3RlZExpc3QgPSBwcm9wcy5uZXN0ZWRJdGVtcy5tYXAoZWwgPT4ge1xyXG4gICAgICAgIGNvbnN0IHBhdGhOYW1lID0gZWwuc2xpY2UoKS5yZXBsYWNlQWxsKFwiIFwiLCBcIi1cIilcclxuICAgICAgICBsZXQgcm91dGUgPSBwYXRoTmFtZSA9PSBcIndlbGNvbWVcIj8gXCJcIiA6IHBhdGhOYW1lXHJcbiAgICAgICAgcmV0dXJuIDxMaW5rIFxyXG4gICAgICAgICAgICAgICAga2V5PXtwcm9wcy5uZXN0ZWRJdGVtcy5pbmRleE9mKGVsKX0gXHJcbiAgICAgICAgICAgICAgICBocmVmPXtgL2RvY3MvJHtyb3V0ZX1gfT5cclxuICAgICAgICAgICAgICAgICAgICA8YT48VG9wRWxlbWVudCB0aXRsZT17ZWx9Lz48L2E+XHJcbiAgICAgICAgICAgICAgICA8L0xpbms+XHJcbiAgICB9KVxyXG5cclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPGRpdiBjbGFzc05hbWU9XCJmbGV4IGZsZXgtY29sIGp1c3RpZnktYmV0d2VlbiBmbGV4LTEgbWwtMTBcIj5cclxuICAgICAgICAgICAge25lc3RlZExpc3R9XHJcbiAgICAgICAgPC9kaXY+XHJcbiAgICApXHJcbn1cclxuXHJcbmZ1bmN0aW9uIE1lbnVFbGVtZW50KHByb3BzKSB7XHJcbiAgICBjb25zdCBbaXNIaWRkZW4sIHNldEhpZGRlbl0gPSB1c2VTdGF0ZShmYWxzZSlcclxuXHJcbiAgICBjb25zdCBoYW5kbGVDbGljayA9IGZ1bmN0aW9uKCkge1xyXG4gICAgICAgIHNldEhpZGRlbighaXNIaWRkZW4pXHJcbiAgICB9IFxyXG5cclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPFJlYWN0LkZyYWdtZW50PlxyXG4gICAgICAgICAgICA8VG9wRWxlbWVudCB0aXRsZT17cHJvcHMudGl0bGV9IG9uQ2xpY2s9e2hhbmRsZUNsaWNrfSAvPlxyXG4gICAgICAgICAgICA8TmVzdGVkTWVudSBoaWRkZW49e2lzSGlkZGVufSBuZXN0ZWRJdGVtcz17cHJvcHMubmVzdGVkSXRlbXN9IC8+XHJcbiAgICAgICAgPC9SZWFjdC5GcmFnbWVudD5cclxuICAgICAgICBcclxuICAgIClcclxufVxyXG4vLy0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLVxyXG5cclxuZXhwb3J0IGRlZmF1bHQgZnVuY3Rpb24gU2lkZUJhcihwcm9wcykge1xyXG4gICAgY29uc3QgeyBtYWluTWVudSwgc2VlQWxzbyB9ID0gc2lkZUJhck1lbnVcclxuXHJcbiAgICBjb25zdCBtZW51TGlzdCA9IG1haW5NZW51Lm1hcChlbCA9PiB7XHJcbiAgICAgICAgcmV0dXJuIChcclxuICAgICAgICAgICAgPE1lbnVFbGVtZW50XHJcbiAgICAgICAgICAgICAgICBrZXk9e21haW5NZW51LmluZGV4T2YoZWwpfVxyXG4gICAgICAgICAgICAgICAgdGl0bGU9e2VsLnRpdGxlfVxyXG4gICAgICAgICAgICAgICAgbmVzdGVkSXRlbXM9e2VsLm5lc3RlZEl0ZW1zfVxyXG4gICAgICAgICAgICAgLz5cclxuICAgICAgICApXHJcbiAgICB9KVxyXG5cclxuICAgIHJldHVybiAoXHJcbiAgICAgICAgPGRpdiBjbGFzc05hbWU9XCJmbGV4IGZsZXgtY29sIHctNjQgaC1zY3JlZW4gcHgtNCBweS04IGJnLXdoaXRlIGJvcmRlci1yXCI+XHJcbiAgICAgICAgICAgIDxkaXY+XHJcbiAgICAgICAgICAgICAgIDxJemUgLz4gXHJcbiAgICAgICAgICAgIDwvZGl2PlxyXG5cclxuICAgICAgICAgICAgPGRpdiBjbGFzc05hbWU9XCJmbGV4IGZsZXgtY29sIGp1c3RpZnktYmV0d2VlbiBmbGV4LTFcIj5cclxuICAgICAgICAgICAgICAgIDxuYXY+XHJcbiAgICAgICAgICAgICAgICAgICAge21lbnVMaXN0fVxyXG4gICAgICAgICAgICAgICAgICAgIDxociBjbGFzc05hbWU9XCJteS02IGJvcmRlci1ncmF5LTIwMFwiIC8+XHJcbiAgICAgICAgICAgICAgICAgICAgPFRvcEVsZW1lbnRcclxuICAgICAgICAgICAgICAgICAgICAgICAgdGl0bGU9e3NlZUFsc28udGl0bGV9XHJcbiAgICAgICAgICAgICAgICAgICAgLz5cclxuICAgICAgICAgICAgICAgIDwvbmF2PlxyXG4gICAgICAgICAgICA8L2Rpdj5cclxuICAgICAgICA8L2Rpdj5cclxuICAgIClcclxufSJdLCJuYW1lcyI6WyJJemUiLCJSZWFjdCIsInVzZVN0YXRlIiwiTGluayIsInNpZGVCYXJNZW51IiwiVG9wRWxlbWVudCIsInByb3BzIiwiZGl2IiwiaWQiLCJ0aXRsZSIsImNsYXNzTmFtZSIsIm9uQ2xpY2siLCJzcGFuIiwiTmVzdGVkTWVudSIsImhpZGRlbiIsIm5lc3RlZExpc3QiLCJuZXN0ZWRJdGVtcyIsIm1hcCIsImVsIiwicGF0aE5hbWUiLCJzbGljZSIsInJlcGxhY2VBbGwiLCJyb3V0ZSIsImhyZWYiLCJhIiwiaW5kZXhPZiIsIk1lbnVFbGVtZW50IiwiaXNIaWRkZW4iLCJzZXRIaWRkZW4iLCJoYW5kbGVDbGljayIsIkZyYWdtZW50IiwiU2lkZUJhciIsIm1haW5NZW51Iiwic2VlQWxzbyIsIm1lbnVMaXN0IiwibmF2IiwiaHIiXSwic291cmNlUm9vdCI6IiJ9\n//# sourceURL=webpack-internal:///./components/sideBar.js\n");

/***/ })

});