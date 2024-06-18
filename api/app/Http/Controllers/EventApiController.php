<?php

namespace App\Http\Controllers;

use App\Services\EventService;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class EventApiController
{
    private $_service;

    public function __construct(EventService $service)
    {
        $this->_service = $service;
    }

    public function all(Request $request, $isAdmin)
    {
        $language = $request->get('lang');
        $search = $request->query('search', '');
        $perPage = $request->query('per_page', 10);
        $sortBy = $request->query('sort_by', '');
        $sortDirection = $request->query('sort_direction', 'desc');

        return $this->_service->all($language, $perPage, $search, $sortBy, $sortDirection, $isAdmin);
    }


    public function allPublic (Request $request)
    {
        return response()->json($this->all($request, false), Response::HTTP_OK);
    }

    public function allAdmin (Request $request)
    {
        return response()->json($this->all($request, true), Response::HTTP_OK);
    }

    public function homePageEvents(Request $request)
    {
        $language = $request->get('lang');

        return response()->json($this->_service->getHomePageEvents($language), Response::HTTP_OK);
    }

    public function getByIdAdmin(Request $request, int $id)
    {
        $language = $request->get('lang');
        return response()->json($this->_service->getById($id, $language, true), Response::HTTP_OK);
    }

    public function add(Request $request)
    {
        $data = $request->all();
        $event = $this->_service->add($data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $event];
    }

    public function update(Request $request, int $id)
    {
        $data = $request->all();
        $event = $this->_service->update($id, $data);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["data" => $event];
    }

    public function delete(int $id)
    {
        $deletedEvent = $this->_service->delete($id);

        if ($this->_service->hasErrors()) {
            return [
                "errors" => $this->_service->getErrors()
            ];
        }

        return ["deleted" => $deletedEvent];
    }
}
