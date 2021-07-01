package main

import (
"fmt"
"github.com/labstack/echo"
_ "github.com/labstack/echo/middleware"
"net/http"
)

func main() {

	e := echo.New();
	e.GET("/", func(c echo.Context) error {
		fmt.Println("핸들러")
		return c.String(http.StatusOK, "hello Wolrd")
	})

	a := e.Group("/a", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("a그룹 미들웨어 시작")
			err := next(c)
			fmt.Println("a그룹 미들웨어 종료")
			return err
		}
	})

	a.GET("/b", func(c echo.Context) error {
		fmt.Println("핸들러 a")
		return c.String(http.StatusOK, "bb")
	})

	a.GET("/b/b", func(c echo.Context) error {
		fmt.Println("핸들러 aa")
		return c.String(http.StatusOK, "bbbb")
	}, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Route 미들웨어 시작")
			err := next(c)
			fmt.Println("Route 미들웨어 종료")
			return err
		}
	})

	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Pre 1번 미들웨어 시작")
			err := next(c)
			fmt.Println("Pre 1번 미들웨어 종료")
			return err
		}
	})

	e.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("Pre 2번 미들웨어 시작")
			err := next(c)
			fmt.Println("Pre 2번 미들웨어 종료")
			return err
		}
	})

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			fmt.Println("미들웨어 시작")
			fmt.Println(c.Path())
			err := next(c)
			fmt.Println("미들웨어 종료")
			return err
		}
	})

	e.Start(":8080")

}