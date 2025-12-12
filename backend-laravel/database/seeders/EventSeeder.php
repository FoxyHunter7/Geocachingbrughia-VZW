<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use App\Models\Event;
use Illuminate\Support\Facades\DB;

class EventSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        DB::table('events')->delete();

        $data = ReadFromCsv::getDataFromCsv(storage_path('app/data/csv/events.csv'), ';');

        $model = new Event();
        foreach ($data as $row) {
            $model->create($row);
        }
    }
}
