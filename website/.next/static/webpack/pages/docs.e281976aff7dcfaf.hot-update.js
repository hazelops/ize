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

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"default\": function() { return /* binding */ SideBar; }\n/* harmony export */ });\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! react/jsx-runtime */ \"./node_modules/react/jsx-runtime.js\");\n/* harmony import */ var react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var _ize__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! ./ize */ \"./components/ize.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2__ = __webpack_require__(/*! react */ \"./node_modules/react/index.js\");\n/* harmony import */ var react__WEBPACK_IMPORTED_MODULE_2___default = /*#__PURE__*/__webpack_require__.n(react__WEBPACK_IMPORTED_MODULE_2__);\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3__ = __webpack_require__(/*! next/link */ \"./node_modules/next/link.js\");\n/* harmony import */ var next_link__WEBPACK_IMPORTED_MODULE_3___default = /*#__PURE__*/__webpack_require__.n(next_link__WEBPACK_IMPORTED_MODULE_3__);\n/* harmony import */ var _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__ = __webpack_require__(/*! ../utilities/sideBarMenu */ \"./utilities/sideBarMenu.js\");\n/* module decorator */ module = __webpack_require__.hmd(module);\n\n\n\n\n\nvar _s = $RefreshSig$();\nfunction TopElement(props) {\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        id: props.title,\n        className: \"flex items-center px-4 py-2 mt-5 text-gray-600 rounded-md hover:bg-gray-200 transition-colors duration-300 transform cursor-pointer\",\n        onClick: props.onClick,\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 8,\n            columnNumber: 9\n        },\n        __self: this,\n        children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"span\", {\n            className: \"mx-4 font-medium capitalize\",\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 9,\n                columnNumber: 13\n            },\n            __self: this,\n            children: props.title\n        })\n    }));\n}\n_c = TopElement;\nfunction NestedMenu(props) {\n    var _this = this;\n    if (props.hidden) {\n        return null;\n    }\n    var nestedList = props.nestedItems.map(function(el) {\n        var pathName = el.slice().replaceAll(\" \", \"-\");\n        var route = pathName == \"welcome\" ? \"\" : pathName;\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)((next_link__WEBPACK_IMPORTED_MODULE_3___default()), {\n            href: \"/docs/\".concat(route),\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 22,\n                columnNumber: 16\n            },\n            __self: _this,\n            children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                title: el,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 25,\n                    columnNumber: 21\n                },\n                __self: _this\n            })\n        }, props.nestedItems.indexOf(el)));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n        className: \"flex flex-col justify-between flex-1 ml-10\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 31,\n            columnNumber: 9\n        },\n        __self: this,\n        children: nestedList\n    }));\n}\n_c1 = NestedMenu;\nfunction MenuElement(props) {\n    _s();\n    var ref = (0,react__WEBPACK_IMPORTED_MODULE_2__.useState)(false), isHidden = ref[0], setHidden = ref[1];\n    var handleClick = function handleClick() {\n        setHidden(!isHidden);\n    };\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)((react__WEBPACK_IMPORTED_MODULE_2___default().Fragment), {\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 45,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(TopElement, {\n                title: props.title,\n                onClick: handleClick,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 46,\n                    columnNumber: 13\n                },\n                __self: this\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(NestedMenu, {\n                hidden: isHidden,\n                nestedItems: props.nestedItems,\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 47,\n                    columnNumber: 13\n                },\n                __self: this\n            })\n        ]\n    }));\n}\n_s(MenuElement, \"Hdw5EO+DplCNBEJcNuH8tsP7WZ4=\");\n_c2 = MenuElement;\n//------------------------------------------------------------------------------------------------------------------\nfunction SideBar(props) {\n    var _this = this;\n    var mainMenu = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.mainMenu, seeAlso = _utilities_sideBarMenu__WEBPACK_IMPORTED_MODULE_4__.sideBarMenu.seeAlso;\n    var menuList = mainMenu.map(function(el) {\n        return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(MenuElement, {\n            title: el.title,\n            nestedItems: el.nestedItems,\n            __source: {\n                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                lineNumber: 59,\n                columnNumber: 13\n            },\n            __self: _this\n        }, mainMenu.indexOf(el)));\n    });\n    return(/*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"div\", {\n        className: \"flex flex-col w-64 h-screen px-4 py-8 bg-white border-r\",\n        __source: {\n            fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n            lineNumber: 68,\n            columnNumber: 9\n        },\n        __self: this,\n        children: [\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 69,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(_ize__WEBPACK_IMPORTED_MODULE_1__[\"default\"], {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 70,\n                        columnNumber: 16\n                    },\n                    __self: this\n                })\n            }),\n            /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"div\", {\n                className: \"flex flex-col justify-between flex-1\",\n                __source: {\n                    fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                    lineNumber: 73,\n                    columnNumber: 13\n                },\n                __self: this,\n                children: /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsxs)(\"nav\", {\n                    __source: {\n                        fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                        lineNumber: 74,\n                        columnNumber: 17\n                    },\n                    __self: this,\n                    children: [\n                        menuList,\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(\"hr\", {\n                            className: \"my-6 border-gray-200\",\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 76,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        }),\n                        /*#__PURE__*/ (0,react_jsx_runtime__WEBPACK_IMPORTED_MODULE_0__.jsx)(MenuElement, {\n                            title: seeAlso.title,\n                            nestedItems: props.generatedDocs,\n                            __source: {\n                                fileName: \"C:\\\\Users\\\\elect\\\\Desktop\\\\ize\\\\website\\\\components\\\\sideBar.js\",\n                                lineNumber: 77,\n                                columnNumber: 21\n                            },\n                            __self: this\n                        })\n                    ]\n                })\n            })\n        ]\n    }));\n};\n_c3 = SideBar;\nvar _c, _c1, _c2, _c3;\n$RefreshReg$(_c, \"TopElement\");\n$RefreshReg$(_c1, \"NestedMenu\");\n$RefreshReg$(_c2, \"MenuElement\");\n$RefreshReg$(_c3, \"SideBar\");\n\n\n;\n    var _a, _b;\n    // Legacy CSS implementations will `eval` browser code in a Node.js context\n    // to extract CSS. For backwards compatibility, we need to check we're in a\n    // browser context before continuing.\n    if (typeof self !== 'undefined' &&\n        // AMP / No-JS mode does not inject these helpers:\n        '$RefreshHelpers$' in self) {\n        var currentExports = module.__proto__.exports;\n        var prevExports = (_b = (_a = module.hot.data) === null || _a === void 0 ? void 0 : _a.prevExports) !== null && _b !== void 0 ? _b : null;\n        // This cannot happen in MainTemplate because the exports mismatch between\n        // templating and execution.\n        self.$RefreshHelpers$.registerExportsForReactRefresh(currentExports, module.id);\n        // A module can be accepted automatically based on its exports, e.g. when\n        // it is a Refresh Boundary.\n        if (self.$RefreshHelpers$.isReactRefreshBoundary(currentExports)) {\n            // Save the previous exports on update so we can compare the boundary\n            // signatures.\n            module.hot.dispose(function (data) {\n                data.prevExports = currentExports;\n            });\n            // Unconditionally accept an update to this module, we'll check if it's\n            // still a Refresh Boundary later.\n            module.hot.accept();\n            // This field is set when the previous version of this module was a\n            // Refresh Boundary, letting us know we need to check for invalidation or\n            // enqueue an update.\n            if (prevExports !== null) {\n                // A boundary can become ineligible if its exports are incompatible\n                // with the previous exports.\n                //\n                // For example, if you add/remove/change exports, we'll want to\n                // re-execute the importing modules, and force those components to\n                // re-render. Similarly, if you convert a class component to a\n                // function, we want to invalidate the boundary.\n                if (self.$RefreshHelpers$.shouldInvalidateReactRefreshBoundary(prevExports, currentExports)) {\n                    module.hot.invalidate();\n                }\n                else {\n                    self.$RefreshHelpers$.scheduleUpdate();\n                }\n            }\n        }\n        else {\n            // Since we just executed the code for the module, it's possible that the\n            // new exports made it ineligible for being a boundary.\n            // We only care about the case when we were _previously_ a boundary,\n            // because we already accepted this update (accidental side effect).\n            var isNoLongerABoundary = prevExports !== null;\n            if (isNoLongerABoundary) {\n                module.hot.invalidate();\n            }\n        }\n    }\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,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Z0JBQUNHLFNBQVMsRUFBQyxDQUFzQzs7Ozs7OztnR0FDaER3QixDQUFHOzs7Ozs7Ozt3QkFDQ0QsUUFBUTs2RkFDUkUsQ0FBRTs0QkFBQ3pCLFNBQVMsRUFBQyxDQUFzQjs7Ozs7Ozs7NkZBQ25DZSxXQUFXOzRCQUNSaEIsS0FBSyxFQUFFdUIsT0FBTyxDQUFDdkIsS0FBSzs0QkFDcEJPLFdBQVcsRUFBRVYsS0FBSyxDQUFDOEIsYUFBYTs7Ozs7Ozs7Ozs7OztBQU14RCxDQUFDO01BL0J1Qk4sT0FBTyIsInNvdXJjZXMiOlsid2VicGFjazovL19OX0UvLi9jb21wb25lbnRzL3NpZGVCYXIuanM/ZmRkYSJdLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgSXplIGZyb20gXCIuL2l6ZVwiXHJcbmltcG9ydCBSZWFjdCwgeyB1c2VTdGF0ZSB9IGZyb20gJ3JlYWN0J1xyXG5pbXBvcnQgTGluayBmcm9tICduZXh0L2xpbmsnXHJcbmltcG9ydCB7IHNpZGVCYXJNZW51IH0gZnJvbSAnLi4vdXRpbGl0aWVzL3NpZGVCYXJNZW51J1xyXG5cclxuZnVuY3Rpb24gVG9wRWxlbWVudChwcm9wcykge1xyXG4gICAgcmV0dXJuIChcclxuICAgICAgICA8ZGl2IGlkPXtwcm9wcy50aXRsZX0gY2xhc3NOYW1lPVwiZmxleCBpdGVtcy1jZW50ZXIgcHgtNCBweS0yIG10LTUgdGV4dC1ncmF5LTYwMCByb3VuZGVkLW1kIGhvdmVyOmJnLWdyYXktMjAwIHRyYW5zaXRpb24tY29sb3JzIGR1cmF0aW9uLTMwMCB0cmFuc2Zvcm0gY3Vyc29yLXBvaW50ZXJcIiBvbkNsaWNrPXtwcm9wcy5vbkNsaWNrfT5cclxuICAgICAgICAgICAgPHNwYW4gY2xhc3NOYW1lPVwibXgtNCBmb250LW1lZGl1bSBjYXBpdGFsaXplXCI+e3Byb3BzLnRpdGxlfTwvc3Bhbj5cclxuICAgICAgICA8L2Rpdj5cclxuICAgIClcclxufVxyXG5cclxuZnVuY3Rpb24gTmVzdGVkTWVudShwcm9wcykge1xyXG4gICAgaWYgKHByb3BzLmhpZGRlbikge1xyXG4gICAgICAgIHJldHVybiBudWxsXHJcbiAgICB9XHJcblxyXG4gICAgY29uc3QgbmVzdGVkTGlzdCA9IHByb3BzLm5lc3RlZEl0ZW1zLm1hcChlbCA9PiB7XHJcbiAgICAgICAgY29uc3QgcGF0aE5hbWUgPSBlbC5zbGljZSgpLnJlcGxhY2VBbGwoXCIgXCIsIFwiLVwiKVxyXG4gICAgICAgIGxldCByb3V0ZSA9IHBhdGhOYW1lID09IFwid2VsY29tZVwiPyBcIlwiIDogcGF0aE5hbWVcclxuICAgICAgICByZXR1cm4gPExpbmsga2V5PXtwcm9wcy5uZXN0ZWRJdGVtcy5pbmRleE9mKGVsKX0gXHJcbiAgICAgICAgICAgIGhyZWY9e2AvZG9jcy8ke3JvdXRlfWB9XHJcbiAgICAgICAgICAgID5cclxuICAgICAgICAgICAgICAgICAgICA8VG9wRWxlbWVudCB0aXRsZT17ZWx9Lz5cclxuICAgICAgICAgICAgICAgIDwvTGluaz5cclxuICAgICAgICBcclxuICAgIH0pXHJcblxyXG4gICAgcmV0dXJuIChcclxuICAgICAgICA8ZGl2IGNsYXNzTmFtZT1cImZsZXggZmxleC1jb2wganVzdGlmeS1iZXR3ZWVuIGZsZXgtMSBtbC0xMFwiPlxyXG4gICAgICAgICAgICB7bmVzdGVkTGlzdH1cclxuICAgICAgICA8L2Rpdj5cclxuICAgIClcclxufVxyXG5cclxuZnVuY3Rpb24gTWVudUVsZW1lbnQocHJvcHMpIHtcclxuICAgIGNvbnN0IFtpc0hpZGRlbiwgc2V0SGlkZGVuXSA9IHVzZVN0YXRlKGZhbHNlKVxyXG5cclxuICAgIGNvbnN0IGhhbmRsZUNsaWNrID0gZnVuY3Rpb24oKSB7XHJcbiAgICAgICAgc2V0SGlkZGVuKCFpc0hpZGRlbilcclxuICAgIH0gXHJcblxyXG4gICAgcmV0dXJuIChcclxuICAgICAgICA8UmVhY3QuRnJhZ21lbnQ+XHJcbiAgICAgICAgICAgIDxUb3BFbGVtZW50IHRpdGxlPXtwcm9wcy50aXRsZX0gb25DbGljaz17aGFuZGxlQ2xpY2t9IC8+XHJcbiAgICAgICAgICAgIDxOZXN0ZWRNZW51IGhpZGRlbj17aXNIaWRkZW59IG5lc3RlZEl0ZW1zPXtwcm9wcy5uZXN0ZWRJdGVtc30gLz5cclxuICAgICAgICA8L1JlYWN0LkZyYWdtZW50PlxyXG4gICAgICAgIFxyXG4gICAgKVxyXG59XHJcbi8vLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tLS0tXHJcblxyXG5leHBvcnQgZGVmYXVsdCBmdW5jdGlvbiBTaWRlQmFyKHByb3BzKSB7XHJcbiAgICBjb25zdCB7IG1haW5NZW51LCBzZWVBbHNvIH0gPSBzaWRlQmFyTWVudVxyXG5cclxuICAgIGNvbnN0IG1lbnVMaXN0ID0gbWFpbk1lbnUubWFwKGVsID0+IHtcclxuICAgICAgICByZXR1cm4gKFxyXG4gICAgICAgICAgICA8TWVudUVsZW1lbnRcclxuICAgICAgICAgICAgICAgIGtleT17bWFpbk1lbnUuaW5kZXhPZihlbCl9XHJcbiAgICAgICAgICAgICAgICB0aXRsZT17ZWwudGl0bGV9XHJcbiAgICAgICAgICAgICAgICBuZXN0ZWRJdGVtcz17ZWwubmVzdGVkSXRlbXN9XHJcbiAgICAgICAgICAgICAvPlxyXG4gICAgICAgIClcclxuICAgIH0pXHJcblxyXG4gICAgcmV0dXJuIChcclxuICAgICAgICA8ZGl2IGNsYXNzTmFtZT1cImZsZXggZmxleC1jb2wgdy02NCBoLXNjcmVlbiBweC00IHB5LTggYmctd2hpdGUgYm9yZGVyLXJcIj5cclxuICAgICAgICAgICAgPGRpdj5cclxuICAgICAgICAgICAgICAgPEl6ZSAvPiBcclxuICAgICAgICAgICAgPC9kaXY+XHJcblxyXG4gICAgICAgICAgICA8ZGl2IGNsYXNzTmFtZT1cImZsZXggZmxleC1jb2wganVzdGlmeS1iZXR3ZWVuIGZsZXgtMVwiPlxyXG4gICAgICAgICAgICAgICAgPG5hdj5cclxuICAgICAgICAgICAgICAgICAgICB7bWVudUxpc3R9XHJcbiAgICAgICAgICAgICAgICAgICAgPGhyIGNsYXNzTmFtZT1cIm15LTYgYm9yZGVyLWdyYXktMjAwXCIgLz5cclxuICAgICAgICAgICAgICAgICAgICA8TWVudUVsZW1lbnRcclxuICAgICAgICAgICAgICAgICAgICAgICAgdGl0bGU9e3NlZUFsc28udGl0bGV9XHJcbiAgICAgICAgICAgICAgICAgICAgICAgIG5lc3RlZEl0ZW1zPXtwcm9wcy5nZW5lcmF0ZWREb2NzfVxyXG4gICAgICAgICAgICAgICAgICAgIC8+XHJcbiAgICAgICAgICAgICAgICA8L25hdj5cclxuICAgICAgICAgICAgPC9kaXY+XHJcbiAgICAgICAgPC9kaXY+XHJcbiAgICApXHJcbn0iXSwibmFtZXMiOlsiSXplIiwiUmVhY3QiLCJ1c2VTdGF0ZSIsIkxpbmsiLCJzaWRlQmFyTWVudSIsIlRvcEVsZW1lbnQiLCJwcm9wcyIsImRpdiIsImlkIiwidGl0bGUiLCJjbGFzc05hbWUiLCJvbkNsaWNrIiwic3BhbiIsIk5lc3RlZE1lbnUiLCJoaWRkZW4iLCJuZXN0ZWRMaXN0IiwibmVzdGVkSXRlbXMiLCJtYXAiLCJlbCIsInBhdGhOYW1lIiwic2xpY2UiLCJyZXBsYWNlQWxsIiwicm91dGUiLCJocmVmIiwiaW5kZXhPZiIsIk1lbnVFbGVtZW50IiwiaXNIaWRkZW4iLCJzZXRIaWRkZW4iLCJoYW5kbGVDbGljayIsIkZyYWdtZW50IiwiU2lkZUJhciIsIm1haW5NZW51Iiwic2VlQWxzbyIsIm1lbnVMaXN0IiwibmF2IiwiaHIiLCJnZW5lcmF0ZWREb2NzIl0sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///./components/sideBar.js\n");

/***/ })

});